package fixtures

import (
	"context"
	"log"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var GooseLocalContainer testcontainers.Container

func GooseLocalInit() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: filepath.Join(".", "docker", "goose"),
		},
		WaitingFor: wait.ForHealthCheck(),
	}
	var err error
	GooseContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	//_, _, err = GooseContainer.Exec(ctx, []string{"goose", "migrations", "up"})
	//if err != nil {
	//	log.Println("its here")
	//	log.Fatal(err)
	//}
}

func GooseLocalDie() {
	if err := GooseContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}
