package fixtures

import (
	"context"
	"log"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var PostgresTestContainer testcontainers.Container

func InitPostgresTest(cfg PostgresConfig) {
	ctx := context.Background()
	var err error

	var port nat.Port

	if PostgresContainer != nil {
		port, err = PostgresContainer.MappedPort(ctx, "5432/tcp")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		port, err = PostgresGenericContainer.MappedPort(ctx, "5432")
		if err != nil {
			log.Fatal(err)
		}
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../docker/pg_test/",
			Dockerfile:    "Dockerfile",
			PrintBuildLog: true,
			KeepImage:     true,
		},
		WaitingFor: wait.ForHealthCheck(),
		Env: map[string]string{
			"POSTGRES_PORT": port.Port(),
			"POSTGRES_HOST": "postgres_gavna",
		},
		Networks: []string{Network.Name},
	}

	PostgresTestContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}
}
