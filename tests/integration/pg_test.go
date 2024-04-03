package integration

import (
	"context"
	"log/slog"
	"testing"

	"testcontainers/tests/fixtures"
	"testcontainers/tests/storage"

	"github.com/stretchr/testify/suite"
)

type PGTests struct {
	suite.Suite
}

func (*PGTests) SetupTest() {
	fixtures.PostgresInit(fixtures.PostgresConfig{
		DbName:   "users",
		User:     "user",
		Password: "password",
	})
}

func (*PGTests) TearDownTest() {
	fixtures.PostgresDie()
}

func TestPGTests(t *testing.T) {
	suite.Run(t, new(PGTests))
}

func (*PGTests) TestPostgresTestContainer() {
	storage.SaveUser(context.Background())
	user := storage.GetUser(context.Background())

	if user.Age == 27 {
		slog.Info("success")
	} else {
		slog.Info("fail")
	}
}
