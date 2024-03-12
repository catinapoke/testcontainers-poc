package integration

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"testcontainers/tests/fixtures"
)

func SaveUser(ctx context.Context) {
	query := fmt.Sprintf(`INSERT INTO users (name, age)
						VALUES($1, $2);`)

	container := fixtures.PostgresContainer
	dbURL, err := container.ConnectionString(ctx)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(ctx, query, "samvel", 27)
	if err != nil {
		log.Println(err)
	}
}
