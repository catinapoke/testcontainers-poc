package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	prod, err := InitNativeKafkaProducer("hu", brokers[0], "0", 10)
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

	cons, err := InitNativeKafkaConsumer("hu", brokers[0], "10000")
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

	assert.Equal(&testing.T{}, resString, key)
}

func InitNativeKafkaProducer(
	clientID string,
	brokers string,
	//login string,
	//pass string,
	//needAuth bool,
	acks string,
	bufMaxMsg int,
) (*kafka.Producer, error) {
	// see details https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	cfg := kafka.ConfigMap{
		"bootstrap.servers":            brokers,
		"client.id":                    clientID,
		"acks":                         acks,
		"queue.buffering.max.messages": bufMaxMsg,
		"go.delivery.reports":          false,
		// "debug":                        "all", // "secury,broker"
	}

	//if needAuth {
	//	cfg["security.protocol"] = "sasl_plaintext"
	//	cfg["sasl.mechanisms"] = "PLAIN"
	//	cfg["sasl.username"] = login
	//	cfg["sasl.password"] = pass
	//}

	p, err := kafka.NewProducer(&cfg)
	if err != nil {
		slog.Error("new producer", err)
		return nil, err
	}

	slog.Info(fmt.Sprintf("kafka producer %s created", clientID))

	return p, nil
}

func InitNativeKafkaConsumer(
	clientID string,
	brokers string,
	group string,
	// login string,
	// pass string,
	// needAuth bool,
) (*kafka.Consumer, error) {
	// see details https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	config := kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          group,
		"client.id":         clientID,
		"auto.offset.reset": "earliest",
		// "enable.auto.commit": "false",
		// "enable.auto.commit":       "true",
		"auto.commit.interval.ms": 3000,
		// "enable.auto.offset.store": "true",
		// "debug":                        "all", // "secury,broker"
	}

	//if needAuth {
	//	config["security.protocol"] = "sasl_plaintext"
	//	config["sasl.mechanisms"] = "PLAIN"
	//	config["sasl.username"] = login
	//	config["sasl.password"] = pass
	//}

	c, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, errors.Wrapf(err, "create kafka consumer")
	}

	slog.Info(fmt.Sprintf("kafka consumer %s created", clientID))

	return c, nil
}
