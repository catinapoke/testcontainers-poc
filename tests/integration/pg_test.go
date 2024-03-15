package integration

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/suite"
	"testcontainers/tests/fixtures"
	"testcontainers/tests/storage"
)

type IntegrationTests struct {
	suite.Suite
}

func (*IntegrationTests) SetupTest() {
	fixtures.PostgresInit()
}

func (*IntegrationTests) TearDownTest() {
	fixtures.PostgresDie()
}

func TestIntegrationTests(t *testing.T) {
	suite.Run(t, new(IntegrationTests))
}

func (*IntegrationTests) TestPostgresTestContainer() {
	storage.SaveUser(context.Background())
	user := storage.GetUser(context.Background())

	if user.Age == 27 {
		slog.Info("success")
	} else {
		slog.Info("fail")
	}
}
