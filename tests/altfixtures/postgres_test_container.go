package altfixtures

import (
	"log"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var PostgresTestContainer *dockertest.Resource

func RunPostgresTestContainer(cfg PostgresConfig, pool *dockertest.Pool) {
	var err error

	PostgresTestContainer, err := pool.BuildAndRunWithOptions("../docker/pg_test/Dockerfile", &dockertest.RunOptions{
		Env: []string{
			"POSTGRES_URL=" + GetPostgresURL(cfg),
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = PostgresTestContainer.Expire(10)
	if err != nil {
		log.Printf("Could not expire PostgresTestContainer: %s", err)
		defer PostgresTestContainer.Close()
	}
}
