package fixtures

import (
	"context"
	"log"
	"testcontainers/tests/fixtures/kraft"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/network"
)

var KafkaContainer *kraft.KafkaContainer

func KafkaInit() {
	ctx := context.Background()

	var err error
	KafkaContainer, err = kraft.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.6.1"),
		network.WithNetwork([]string{"kafka"}, Network),
		testcontainers.WithEnv(map[string]string{
			"KAFKA_LISTENERS":                      "PLAINTEXT://0.0.0.0:9093,BROKER://0.0.0.0:9092,CONTROLLER://kafka:9094",
			"KAFKA_REST_BOOTSTRAP_SERVERS":         "PLAINTEXT://0.0.0.0:9093,BROKER://0.0.0.0:9092,CONTROLLER://kafka:9094",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP": "BROKER:PLAINTEXT,PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
			"KAFKA_ADVERTISED_LISTENERS":           "PLAINTEXT://0.0.0.0:9093,BROKER://kafka:9092,CONTROLLER://kafka:9094",
			"KAFKA_INTER_BROKER_LISTENER_NAME":     "BROKER",
			"KAFKA_CONTROLLER_LISTENER_NAMES":      "CONTROLLER",
		}),
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
