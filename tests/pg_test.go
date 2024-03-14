package tests

import (
	"context"
	"log/slog"
	"testing"

	"testcontainers/tests/fixtures"
	"testcontainers/tests/integration"
)

func TestMain(m *testing.M) {
	fixtures.PostgresInit()
	defer fixtures.PostgresDie()

	m.Run()
}

func TestPostgresTestContainer(t *testing.T) {
	integration.SaveUser(context.Background())
	user := integration.GetUser(context.Background())

	if user.Age == 27 {
		slog.Info("success")
	} else {
		slog.Info("fail")
	}
}
