package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var GooseContainer testcontainers.Container

func PostgresGooseInit(cfg PostgresConfig) {
	ctx := context.Background()
	var err error

	//dbstring, err := PostgresContainer.ConnectionString(ctx)
	//if err != nil {
	//	log.Fatal(err)
	//}

	containerPort, err := PostgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatal(err)
	}

	host, err := PostgresContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dbstring := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, containerPort.Port(), "user", "password", "users")
	n := Network.Name
	req := testcontainers.ContainerRequest{
		Image:      "ghcr.io/kukymbr/goose-docker:3.19.2",
		WaitingFor: wait.ForHealthCheck(),
		Env: map[string]string{
			"GOOSE_COMMAND":  "up",
			"GOOSE_DRIVER":   "postgres",
			"GOOSE_DBSTRING": dbstring,
		},
		Networks: []string{n},
		NetworkAliases: map[string][]string{
			Network.Name: {"alias1", "alias2", "alias3"},
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

	GooseContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = GooseContainer.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
