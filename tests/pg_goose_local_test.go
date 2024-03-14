package tests

//
//import (
//	"log/slog"
//	"testing"
//
//	"testcontainers/tests/fixtures"
//	"testcontainers/tests/integration"
//)
//
//func TestPostgresWithGooseLocalTestContainer(t *testing.T) {
//	fixtures.PostgresInit()
//	defer fixtures.PostgresDie()
//
//	fixtures.GooseLocalInit()
//	defer fixtures.GooseLocalDie()
//
//	integration.SaveUser(ctx)
//	user := integration.GetUser(ctx)
//
//	if user.Age == 27 {
//		slog.Info("success")
//	} else {
//		slog.Info("fail")
//	}
//}
