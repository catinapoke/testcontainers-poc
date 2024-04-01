package integration

import (
	"context"
	"log"
	"log/slog"
	"testcontainers/tests/altfixtures"
	"testcontainers/tests/storage"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type PGGooseDockerTests struct {
	suite.Suite
	pool *dockertest.Pool
}

func (t *PGGooseDockerTests) SetupTest() {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	t.pool = pool

	cfg := altfixtures.PostgresConfig{
		Name:     "postgres",
		User:     "postgres",
		Password: "postgres",
	}

	// fixtures.NetworkInit()
	altfixtures.PostgresInit(cfg, pool)
	altfixtures.RunPostgresTestContainer(cfg, pool)
}

func (t *PGGooseDockerTests) TearDownTest() {
	altfixtures.PostgresDie(t.pool)
}

func TestDockerTests(t *testing.T) {
	suite.Run(t, new(PGGooseDockerTests))
}

func (p *PGGooseDockerTests) TestPostgresTestContainer() {

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
