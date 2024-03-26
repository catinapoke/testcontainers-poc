package integration

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"testcontainers/tests/fixtures"
	"testcontainers/tests/storage"
)

type PGGooseTests struct {
	suite.Suite
}

func (*PGGooseTests) SetupTest() {
	fixtures.NetworkInit()
	fixtures.PostgresInit(fixtures.PostgresConfig{
		Name:     "users",
		User:     "user",
		Password: "password",
	})
	fixtures.PostgresGooseInit(fixtures.PostgresConfig{
		Name:     "users",
		User:     "user",
		Password: "password",
	})
}

func (*PGGooseTests) TearDownTest() {
	fixtures.PostgresDie()
	fixtures.NetworkDie()
}

func TestPGGooseTests(t *testing.T) {
	suite.Run(t, new(PGGooseTests))
}

func (p *PGGooseTests) TestPostgresTestContainer() {
	//a, err := fixtures.PostgresContainer.ConnectionString(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//slog.Info(a)

	//fixtures.PostgresGooseInit(fixtures.PostgresConfig{
	//	Name:     "users",
	//	User:     "user",
	//	Password: "password",
	//})

	time.Sleep(time.Minute)
	storage.SaveUser(context.Background())
	user := storage.GetUser(context.Background())

	if user.Age == 27 {
		slog.Info("success")
	} else {
		slog.Info("fail")
	}

	p.Suite.Equal(int64(27), int64(user.Age))
}
