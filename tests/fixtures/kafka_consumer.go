package fixtures

import (
	"fmt"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pkg/errors"
)

func InitNativeKafkaConsumer(
	clientID string,
	brokers string,
	group string,
) (*kafka.Consumer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers":       brokers,
		"group.id":                group,
		"client.id":               clientID,
		"auto.offset.reset":       "earliest",
		"auto.commit.interval.ms": 3000,
	}

	c, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "create kafka consumer")
	}

	slog.Info(fmt.Sprintf("kafka consumer %s created", clientID))

	return c, nil
}
