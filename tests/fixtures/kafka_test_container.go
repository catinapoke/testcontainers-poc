package fixtures

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var KafkaTestContainer testcontainers.Container

func InitKafkaTest(brokers string, input string, output string) {
	ctx := context.Background()
	var err error

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../docker/kafka_rw_test/",
			Dockerfile:    "Dockerfile",
			PrintBuildLog: true,
			KeepImage:     true,
		},
		WaitingFor: wait.ForLog("start consuming events"),
		Env: map[string]string{
			"KAFKA_BROKERS":   "kafka:9093",
			"KAFKA_TOPIC_IN":  input,
			"KAFKA_TOPIC_OUT": output,
		},
		Networks: []string{Network.Name},
	}

	KafkaTestContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func KafkaTestDie() {
	if err := KafkaTestContainer.Terminate(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}
}

func CreateTopic(ctx context.Context, brokers []string, topics []string) error {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": strings.Join(brokers, ";")})
	defer admin.Close()

	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	specs := make([]kafka.TopicSpecification, 0, len(topics))

	for _, item := range topics {
		specs = append(specs, kafka.TopicSpecification{
			Topic:         item,
			NumPartitions: 1,
		})
	}

	_, err = admin.CreateTopics(ctx, specs)

	if err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	fmt.Println("succesfully created topics")

	return nil
}
