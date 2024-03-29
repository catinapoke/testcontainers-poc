package fixtures

import (
	"context"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

var PostgresContainer *postgres.PostgresContainer
var PostgresGenericContainer testcontainers.Container

type PostgresConfig struct {
	Name     string
	User     string
	Password string
}

func PostgresInit(cfg PostgresConfig) {
	ctx := context.Background()

	var err error
	PostgresContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-alpine3.19"), //docker.io/postgres:15.2-alpine
		postgres.WithConfigFile("../postgres/my-postgres.conf"),
		// postgres.WithInitScripts("../../migrations/0001_create_users_table.sql"),
		postgres.WithDatabase(cfg.Name),
		postgres.WithUsername(cfg.User),
		postgres.WithPassword(cfg.Password),
		network.WithNetwork([]string{"postgres_gavna"}, Network),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
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
			"POSTGRES_DB":       cfg.Name,
			"POSTGRES_USER":     cfg.User,
			"POSTGRES_PASSWORD": cfg.Password,
			"listen_addresses":  "*",
		},
		ExposedPorts: []string{"5432"},
		Networks:     []string{Network.Name},
		NetworkAliases: map[string][]string{
			Network.Name: {"postgres_gavna"},
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second),
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
