package integration

import (
	"context"
	"log/slog"
	"testing"

	"testcontainers/tests/fixtures"
	"testcontainers/tests/storage"

	"github.com/stretchr/testify/suite"
)

type PGGooseTests struct {
	suite.Suite
}

func (*PGGooseTests) SetupTest() {
	cfg := fixtures.PostgresConfig{
		DbName:   "postgres",
		User:     "postgres",
		Password: "postgres",
		Port:     "5432/tcp",
	}

	fixtures.NetworkInit()
	fixtures.PostgresInit(cfg) //PostgresInit
	// fixtures.InitGooseFromDockerfile(context.TODO(), fixtures.Network, cfg)
	fixtures.InitPostgresGoose(cfg)
	fixtures.InitPostgresTest(cfg)
}

func (*PGGooseTests) TearDownTest() {
	fixtures.PostgresDie()
	fixtures.GooseContainer.Terminate(context.TODO())
	fixtures.NetworkDie()
}

func TestPGGooseTests(t *testing.T) {
	suite.Run(t, new(PGGooseTests))
}

func (p *PGGooseTests) TestPostgresTestContainer() {
	storage.SaveUser(context.Background())
	user := storage.GetUser(context.Background())

	if user.Age == 27 {
		slog.Info("success")
	} else {
		slog.Info("fail")
	}

	p.Suite.Equal(int64(27), int64(user.Age))
}
