package fixtures

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var GooseContainer testcontainers.Container

func PostgresGooseInit(cfg PostgresConfig) {
	ctx := context.Background()

	dbURL, err := PostgresContainer.ConnectionString(ctx)
	if err != nil {
		log.Fatal(err)
	}

	req := testcontainers.ContainerRequest{
		Image:      "ghcr.io/kukymbr/goose-docker:3.19.2",
		WaitingFor: wait.ForExit(),
		Env: map[string]string{
			"GOOSE_COMMAND":  "up",
			"GOOSE_DRIVER":   "postgres",
			"GOOSE_DBSTRING": dbURL,
		},
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

	//var err error
	GooseContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = GooseContainer.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
