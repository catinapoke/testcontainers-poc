package helpers

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

func MakeMsg(topic, key string, message interface{}) kafka.Message {
	headers := []kafka.Header{
		{
			Key: "DateAdd", Value: []byte(time.Now().Format(time.RFC3339Nano)),
		},
		{
			Key: "MessageId", Value: []byte(uuid.NewString()),
		},
	}

	messageJson, _ := json.Marshal(message)

	keyJson, _ := json.Marshal(key)

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
