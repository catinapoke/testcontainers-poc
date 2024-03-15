package integration

import (
	"context"
	"log/slog"
	"testing"

	"testcontainers/tests/fixtures"
	"testcontainers/tests/storage"
)

func TestMain(m *testing.M) {
	fixtures.PostgresInit()
	defer fixtures.PostgresDie()

	m.Run()
}

func TestPostgresTestContainer(t *testing.T) {
	storage.SaveUser(context.Background())
	user := storage.GetUser(context.Background())

	if user.Age == 27 {
		slog.Info("success")
	} else {
		slog.Info("fail")
	}
}
