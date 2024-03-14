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
//func TestPostgresWithGooseTestContainer(t *testing.T) {
//	ctx := fixtures.PostgresInit()
//	defer fixtures.PostgresDie()
//
//	fixtures.GooseInit()
//	defer fixtures.GooseDie()
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
