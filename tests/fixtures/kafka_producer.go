package fixtures

import (
	"fmt"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func InitNativeKafkaProducer(
	clientID string,
	brokers string,
	acks string,
	bufMaxMsg int,
) (*kafka.Producer, error) {
	cfg := kafka.ConfigMap{
		"bootstrap.servers":            brokers,
		"client.id":                    clientID,
		"acks":                         acks,
		"queue.buffering.max.messages": bufMaxMsg,
		"go.delivery.reports":          false,
	}

	p, err := kafka.NewProducer(&cfg)
	if err != nil {
		slog.Error("new producer", err)
		return nil, err
	}

	slog.Info(fmt.Sprintf("kafka producer %s created", clientID))

	return p, nil
}
