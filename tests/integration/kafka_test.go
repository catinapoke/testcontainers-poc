package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"testcontainers/tests/fixtures"
	"testcontainers/tests/helpers"
	"testcontainers/tests/storage"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

const (
	topic   = "jeppa"
	key     = "hui"
	message = "bolshoy"

	topic_in  = "pizza"
	topic_out = "cardboard"
)

type KafkaTests struct {
	suite.Suite
}

func (tests *KafkaTests) SetupSuite() {
	fixtures.NetworkInit()
	fixtures.KafkaInit()
	fixtures.PostgresInit(fixtures.PostgresConfig{
		DbName:   "users",
		User:     "user",
		Password: "password",
	})

	brokers, err := fixtures.KafkaContainer.Brokers(context.TODO())
	if err != nil {
		tests.FailNow("failed to get brokers", err)
	}

	fixtures.CreateTopic(context.TODO(), brokers, []string{topic_in, topic_out})

	// brokers_str := strings.Replace(strings.Join(brokers, ","), "localhost", "172.28.200.12", 1)
	fixtures.InitKafkaTest("kafka:9092", topic_in, topic_out) // strings.Replace(strings.Join(brokers, ","), "localhost", "172.28.200.12", 1)
}

func (*KafkaTests) TearDownSuite() {
	fixtures.KafkaDie()
	fixtures.PostgresDie()
	fixtures.KafkaTestDie()
	fixtures.NetworkDie()
}

func TestKafkaTests(t *testing.T) {
	suite.Run(t, new(KafkaTests))
}

func (k *KafkaTests) TestPostgresTestContainer() {
	storage.SaveUser(context.Background())
	user := storage.GetUser(context.Background())

	k.Suite.Equal(int64(27), user.Age)
}

// переписать на t.Error или подумать, как сделать, что бы тест стал Fail, если будет ошибка
// уверен, что достаточно просто возвращать ошибку
// потому что все функции внутри Suite - становятся тестами при запуске
func (k *KafkaTests) TestKafkaTestContainer() {
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

	//msg := helpers.MakeMsg(topic, key, message)
	msg := helpers.MakeMsg(topic, key, storage.User{Name: "Samik", Age: 27})

	prod.Produce(&msg, nil)

	// думаю, что это можно вынести куда то, где мы "на фоне" прочитаем из кафки
	// и запишем в бд
	// а в этом коде чуть ниже мы просто проверим, есть ли запись в БД
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

	var fromKafka storage.User
	err = json.Unmarshal(res.Value, &fromKafka)
	if err != nil {
		log.Fatal("unmarshall error", err)
	}

	storage.SaveUserWithParams(ctx, fromKafka.Name, fromKafka.Age)
	fromDB := storage.GetUserByName(ctx, fromKafka.Name)

	k.Suite.Equal(fromKafka, fromDB)
}

func (k *KafkaTests) TestKafkaConnectivity() {
	ctx := context.Background()

	// Init
	brokers, err := fixtures.KafkaContainer.Brokers(ctx)
	if err != nil {
		k.FailNow("Get Kafka brokers", err)
	}

	producer, err := fixtures.InitNativeKafkaProducer("connectivity-producer", strings.Join(brokers, ","), "1", 10)
	if err != nil {
		k.FailNow("InitNativeKafkaProducer", err)
	}
	defer producer.Close()

	consumer, err := fixtures.InitNativeKafkaConsumer("connectivity-consumer", strings.Join(brokers, ","), "10000")
	if err != nil {
		k.FailNow("InitNativeKafkaConsumer", err)
	}
	defer consumer.Close()

	// Act
	text_msg := "test-input-external"
	msg := helpers.MakeMsg(topic_in, key, text_msg)

	err = producer.Produce(&msg, nil)
	if ok := k.NoError(err, "failed to produce message from external"); !ok {
		return
	}

	err = consumer.Subscribe(topic_out, nil)
	if ok := k.NoError(err, "failed to subscribe to topic from external"); !ok {
		return
	}

	result, err := consumer.ReadMessage(time.Second * 3)
	if ok := k.NoError(err, "failed to read message from external"); !ok {
		WriteLogs(fixtures.KafkaTestContainer)
		return
	}

	// Assert
	k.Contains(string(result.Value), text_msg, "got wrong string")
}

func WriteLogs(container testcontainers.Container) {
	reader, err1 := container.Logs(context.Background())
	defer reader.Close()

	name, err2 := container.Name(context.TODO())
	if err2 != nil {
		name = "default"
	}

	if err1 == nil {
		data, er := io.ReadAll(reader)
		if er == nil {
			fmt.Printf("container %s logs:\n %s", name, string(data))
		}
	}
}
