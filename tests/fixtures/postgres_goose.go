package fixtures

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var GooseContainer testcontainers.Container

func InitPostgresGoose(cfg PostgresConfig) {
	ctx := context.Background()
	var err error

	dbstring := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", PostgresAlias, cfg.Port.Port(), cfg.User, cfg.Password, cfg.DbName)

	req := testcontainers.ContainerRequest{
		Image:      "ghcr.io/kukymbr/goose-docker:3.19.2",
		WaitingFor: wait.ForLog("successfully migrated database"),
		Env: map[string]string{
			"GOOSE_COMMAND":  "up",
			"GOOSE_DRIVER":   "postgres",
			"GOOSE_DBSTRING": dbstring,
		},
		Networks: []string{Network.Name},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath: "../../migrations",
				// ContainerFile cannot create the parent directory, so we copy the scripts
				// to the root of the container instead. Make sure to create the container directory
				// before you copy a host directory on create.
				ContainerFilePath: "/migrations",
				FileMode:          0o700,
			},
		},
	}

	net, _ := PostgresContainer.Networks(ctx)
	log.Print(net)

	GooseContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func RunGooseStatusManual(cfg PostgresConfig) {
	ctx := context.Background()
	var err error

	containerPort, err := PostgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatal(err)
	}

	db, err := InitPgDb(ctx, "localhost", cfg.User, cfg.Password, containerPort.Port(), cfg.DbName, 5, 10)
	if err != nil {
		log.Fatal(err)
	}

	err = goose.Status(db, "../../migrations")
	if err != nil {
		log.Fatal(err)
	}
}

func InitGooseFromDockerfile(ctx context.Context, network *testcontainers.DockerNetwork, cfg PostgresConfig) { // , postgresAlias string
	containerPort, err := PostgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatal(err)
	}

	dbstring := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "postgres_gavna", containerPort.Port(), cfg.User, cfg.Password, cfg.DbName)

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../docker/goose/",
			Dockerfile:    "Dockerfile",
			PrintBuildLog: true,
			KeepImage:     true,
		},
		Env: map[string]string{
			"GOOSE_COMMAND":       "status",
			"GOOSE_DRIVER":        "postgres",
			"GOOSE_MIGRATION_DIR": "/migrations",
			"GOOSE_DBSTRING":      dbstring,
		},
		WaitingFor: wait.ForHealthCheck(),
		Networks:   []string{network.Name},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath: "../../migrations",
				// ContainerFile cannot create the parent directory, so we copy the scripts
				// to the root of the container instead. Make sure to create the container directory
				// before you copy a host directory on create.
				ContainerFilePath: "/migrations",
				FileMode:          0o700,
			},
		},
	}

	GooseContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          false,
	})

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		err := GooseContainer.Start(ctx)

		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

}

func InitPgDb(
	ctx context.Context, host, user, password, port, dbName string,
	openConn, maxOpenConn int,
) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		dbName,
	)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, errors.Wrapf(err, "connect to db %s", host)
	}

	db.SetMaxIdleConns(openConn)
	db.SetMaxOpenConns(maxOpenConn)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "ping db %s", dbName)
	}

	// logger.Logger().Infof("connect to PG established: [%s].[%s]", host, dbName)

	return db, nil
}
