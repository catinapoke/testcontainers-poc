package fixtures

import (
	"context"
	"fmt"
	"log"
	"testcontainers/tests/helpers"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

var PostgresContainer *postgres.PostgresContainer
var PostgresGenericContainer testcontainers.Container
var PostgresAlias string = "pg_database"
var PostgresDefaultConfig = PostgresConfig{
	DbName:   "postgres",
	User:     "postgres",
	Password: "postgres",
	Port:     "5432/tcp",
}

type PostgresConfig struct {
	DbName   string
	User     string
	Password string
	Port     nat.Port
}

func (p PostgresConfig) urlFromPort(port nat.Port) string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", p.User, p.Password, port.Port(), p.DbName)
}

func (p PostgresConfig) urlFromHostPort(host string, port nat.Port) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", p.User, p.Password, host, port.Port(), p.DbName)
}

func PostgresInit(cfg PostgresConfig) {
	ctx := context.Background()

	var err error
	PostgresContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-alpine3.19"), //docker.io/postgres:15.2-alpine
		postgres.WithConfigFile("../postgres/my-postgres.conf"),
		// postgres.WithInitScripts("../../migrations/0001_create_users_table.sql"),
		postgres.WithDatabase(cfg.DbName),
		postgres.WithUsername(cfg.User),
		postgres.WithPassword(cfg.Password),
		network.WithNetwork([]string{PostgresAlias}, Network),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
			helpers.WaitForSQL{
				UrlFromPort: cfg.urlFromPort,
				Driver:      "pgx",
				Port:        cfg.Port,
			},
		),
	)

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	state, err := PostgresContainer.State(ctx)
	if err != nil {
		log.Fatalf("failed to get container state: %s", err) // nolint:gocritic
	}

	log.Println(state.Running)
}

func PostgresInitGeneric(cfg PostgresConfig) {
	ctx := context.Background()

	var err error
	req := testcontainers.ContainerRequest{
		Image: "postgres:16.2-alpine3.19",
		Env: map[string]string{
			"POSTGRES_DB":       cfg.DbName,
			"POSTGRES_USER":     cfg.User,
			"POSTGRES_PASSWORD": cfg.Password,
			"listen_addresses":  "*",
		},
		ExposedPorts: []string{"5432/tcp"},
		Networks:     []string{Network.Name},
		NetworkAliases: map[string][]string{
			Network.Name: {PostgresAlias},
		},
		WaitingFor: helpers.WaitForSQL{
			UrlFromPort: cfg.urlFromPort,
			Driver:      "pgx",
			Port:        cfg.Port,
		},
		// wait.ForLog("database system is ready to accept connections").
		// 	WithOccurrence(2).
		// 	WithStartupTimeout(5 * time.Second),
	}

	PostgresGenericContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalf("failed to start postgres container: %s", err)
	}

	state, err := PostgresGenericContainer.State(ctx)
	if err != nil {
		log.Fatalf("failed to get container state: %s", err) // nolint:gocritic
	}

	log.Println(state.Running)
}

func PostgresDie() {
	if PostgresContainer != nil {
		if err := PostgresContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}

	if PostgresGenericContainer != nil {
		if err := PostgresGenericContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}
}
