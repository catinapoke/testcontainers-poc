package integration

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testcontainers/tests/fixtures"
)

type KafkaTests struct {
	suite.Suite
}

func (*KafkaTests) SetupTest() {
	fixtures.KafkaInit()
}

func (*KafkaTests) TearDownTest() {
	fixtures.KafkaDie()
}

func TestKafkaTests(t *testing.T) {
	suite.Run(t, new(KafkaTests))
}

func (*KafkaTests) TestKafkaTestContainer() {
	ctx := context.Background()
	brokers, err := fixtures.KafkaContainer.Brokers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	prod, err := fixtures.InitNativeKafkaProducer("hu", brokers[0], "0", 10)
	if err != nil {
		log.Fatal("InitNativeKafkaProducer", err)
	}
	defer prod.Close()

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

	prod.Produce(&msg, nil)

	cons, err := fixtures.InitNativeKafkaConsumer("hu", brokers[0], "10000")
	if err != nil {
		log.Fatal("InitNativeKafkaConsumer", err)
	}
	defer cons.Close()

	err = cons.Subscribe(topic, nil)
	if err != nil {
		log.Fatal("subscribe", err)
	}

	res, err := cons.ReadMessage(time.Second * 100)
	if err != nil {
		log.Fatal("read message", err)
	}

	var resString string
	err = json.Unmarshal(res.Key, &resString)
	if err != nil {
		log.Fatal("unmarshall error", err)
	}

	assert.Equal(&testing.T{}, resString, "balshoy")
}
