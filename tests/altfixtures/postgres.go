package altfixtures

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var PostgresContainer *dockertest.Resource

type PostgresConfig struct {
	Name     string
	User     string
	Password string
}

func PostgresInit(cfg PostgresConfig, pool *dockertest.Pool) {
	var err error

	// pulls an image, creates a container based on it and runs it
	PostgresContainer, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.2-alpine3.19",
		Env: []string{
			"POSTGRES_USER=" + cfg.User,
			"POSTGRES_PASSWORD=" + cfg.Password,
			"POSTGRES_DB=" + cfg.Name,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	url := GetPostgresURL(cfg)

	log.Println("Connecting to database on url: ", url)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 10 * time.Second
	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", url)
		defer db.Close()

		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
}

func PostgresDie(pool *dockertest.Pool) {
	if err := pool.Purge(PostgresContainer); err != nil {
		log.Printf("Could not purge PostgresContainer: %s", err)
	}
}

func GetPostgresURL(cfg PostgresConfig) string {
	hostAndPort := PostgresContainer.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.User, cfg.Password, hostAndPort, cfg.Name)
	return databaseUrl
}
