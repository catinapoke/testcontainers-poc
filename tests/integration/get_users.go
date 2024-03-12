package integration

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"testcontainers/tests/fixtures"
)

func GetUser(ctx context.Context) User {
	query := fmt.Sprintf(`select name, age from users
						where name = $1`)

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

	var user User
	err = conn.QueryRow(ctx, query, "samvel").Scan(&user.Name, &user.Age)
	if err != nil {
		log.Println(err)
	}

	return user
}
