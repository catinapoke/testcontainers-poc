package fixtures

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

var KafkaContainer *kafka.KafkaContainer

func KafkaInit() {
	ctx := context.Background()

	var err error
	KafkaContainer, err = kafka.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
}

func KafkaDie() {
	if err := KafkaContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}
