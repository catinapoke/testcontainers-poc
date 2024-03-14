package fixtures

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
)

var GooseContainer testcontainers.Container

func GooseInit() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "gomicro/goose:latest",
	}
	var err error
	GooseContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = GooseContainer.Exec(ctx, []string{"goose", "migrations", "up"})
	if err != nil {
		log.Println("its here")
		log.Fatal(err)
	}
}

func GooseDie() {
	if err := GooseContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}
