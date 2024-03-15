package fixtures

import (
	"context"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var PostgresContainer *postgres.PostgresContainer

func PostgresInit() {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	var err error
	PostgresContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithInitScripts("../../migrations/0001_create_users_table.sql"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
}

func PostgresDie() {
	if err := PostgresContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}
