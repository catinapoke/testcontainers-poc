package fixtures

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var PostgresTestContainer testcontainers.Container

func InitPostgresTest(cfg PostgresConfig) {
	ctx := context.Background()
	var err error

	// var mappedPort nat.Port

	// if PostgresContainer != nil {
	// 	mappedPort, err = PostgresContainer.MappedPort(ctx, cfg.Port)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	mappedPort, err = PostgresGenericContainer.MappedPort(ctx, cfg.Port)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../docker/pg_test/",
			Dockerfile:    "Dockerfile",
			PrintBuildLog: true,
			KeepImage:     true,
		},
		WaitingFor: wait.ForLog("Succesfully connected"),
		Env: map[string]string{
			// "POSTGRES_PORT": port.Port(),
			// "POSTGRES_HOST": PostgresAlias,
			"POSTGRES_URL": cfg.urlFromHostPort(PostgresAlias, cfg.Port),
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
