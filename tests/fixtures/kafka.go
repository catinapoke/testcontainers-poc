package fixtures

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/network"
)

var KafkaContainer *kafka.KafkaContainer

func KafkaInit() {
	ctx := context.Background()

	var err error
	KafkaContainer, err = kafka.RunContainer(ctx,
		kafka.WithClusterID("test-cluster"),
		testcontainers.WithImage("confluentinc/confluent-local:7.5.0"),
		network.WithNetwork([]string{"kafka"}, Network),
		testcontainers.WithEnv(map[string]string{
			"KAFKA_LISTENERS":                      "INSIDE://0.0.0.0:9095, PLAINTEXT://0.0.0.0:9093,BROKER://0.0.0.0:9092,CONTROLLER://0.0.0.0:9094",
			"KAFKA_REST_BOOTSTRAP_SERVERS":         "INSIDE://0.0.0.0:9095, PLAINTEXT://0.0.0.0:9093,BROKER://0.0.0.0:9092,CONTROLLER://0.0.0.0:9094",
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP": "INSIDE:PLAINTEXT,BROKER:PLAINTEXT,PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT",
			"KAFKA_ADVERTISED_LISTENERS":           "INSIDE://kafka:9095, PLAINTEXT://0.0.0.0:9093",
			// "KAFKA_INTER_BROKER_LISTENER_NAME":     "INSIDE",
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
