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

	port, _ := os.LookupEnv("POSTGRES_PORT")
	host, _ := os.LookupEnv("POSTGRES_HOST")

	db, err := InitPgDb(context.Background(), host, cfg.User, cfg.Password, port, cfg.Name, 5, 10)
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
}

func InitPgDb(
	ctx context.Context, host, user, password, port, dbName string,
	openConn, maxOpenConn int,
) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		url.QueryEscape(user),
		url.QueryEscape(password),
		host,
		port,
		dbName,
	)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, errors.Wrapf(err, "connect to db %s", host)
	}

	db.SetMaxIdleConns(openConn)
	db.SetMaxOpenConns(maxOpenConn)

	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "ping db %s", dbName)
	}

	// logger.Logger().Infof("connect to PG established: [%s].[%s]", host, dbName)

	return db, nil
}
