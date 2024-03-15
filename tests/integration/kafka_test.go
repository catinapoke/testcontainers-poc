package integration

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testcontainers/tests/fixtures"
	"testcontainers/tests/helpers"
)

const (
	topic   = "jeppa"
	key     = "hui"
	message = "bolshoy"
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

	msg := helpers.MakeMsg(topic, key, message)

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

	assert.Equal(&testing.T{}, resString, key)
}
