package fixtures

import (
	"fmt"
	"log"
	"log/slog"
	"time"

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

func Consume(broker string) *kafka.Message {
	cons, err := InitNativeKafkaConsumer("hu", broker, "10000")
	if err != nil {
		log.Fatal("InitNativeKafkaConsumer", err)
	}
	defer cons.Close()

	err = cons.Subscribe("jeppa", nil)
	if err != nil {
		log.Fatal("subscribe", err)
	}

	res, err := cons.ReadMessage(time.Second * 100)
	if err != nil {
		log.Fatal("read message", err)
	}

	return res
}
