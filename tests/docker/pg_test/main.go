package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	type Config struct {
		Name     string
		User     string
		Password string
	}

	cfg := Config{
		Name:     "postgres",
		User:     "postgres",
		Password: "postgres",
	}

	url, _ := os.LookupEnv("POSTGRES_URL")

	if len(url) == 0 {
		port, _ := os.LookupEnv("POSTGRES_PORT")
		host, _ := os.LookupEnv("POSTGRES_HOST")

		url = GetConnectionUrl(host, cfg.User, cfg.Password, port, cfg.Name)
	}

	fmt.Printf("Ready for connect to %s\n", url)

	db, err := InitPgDb(context.Background(), url, 5, 10)
	if err != nil {
		fmt.Printf("failed to connect\n")
		log.Fatal(err)
	}

	db.Close()

	fmt.Print("Succesfully connected\n")
}

func InitPgDb(
	ctx context.Context, url string,
	openConn, maxOpenConn int,
) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, errors.Wrapf(err, "connect to db %s", url)
	}

	db.SetMaxIdleConns(openConn)
	db.SetMaxOpenConns(maxOpenConn)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "ping db %s", url)
	}

	return db, nil
}

func GetConnectionUrl(host, user, password, port, dbName string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		dbName,
	)
}
