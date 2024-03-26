package fixtures

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

var Network *testcontainers.DockerNetwork

func NetworkInit() {
	ctx := context.Background()
	var err error

	Network, err = network.New(ctx, network.WithCheckDuplicate())
	if err != nil {
		log.Fatal(err)
	}
}

func NetworkDie() {
	if err := Network.Remove(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}
