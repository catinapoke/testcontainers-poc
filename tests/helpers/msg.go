package helpers

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

func MakeMsg() kafka.Message {
	headers := []kafka.Header{
		{
			Key: "DateAdd", Value: []byte(time.Now().Format(time.RFC3339Nano)),
		},
		{
			Key: "MessageId", Value: []byte(uuid.NewString()),
		},
	}

	message := "hui"
	messageJson, _ := json.Marshal(message)

	key := "balshoy"
	keyJson, _ := json.Marshal(key)

	topic := "jeppa"

	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value:   messageJson,
		Key:     keyJson,
		Headers: headers,
	}

	return msg
}
