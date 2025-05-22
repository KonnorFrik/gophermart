package model

import (
	"errors"
	db "gophermart/model/postgres"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	// pgx
	// sqlc - https://github.com/sqlc-dev/sqlc
	// https://github.com/doug-martin/goqu
	// migrations
)

var (
	config = &db.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DbName:   os.Getenv("DB_DB"),
		Port:     os.Getenv("DB_PORT"),
		TimeZone: "Asia/Yekaterinburg",
	}

	dbObj *pgx.Conn
    ErrDataBaseNotConnected = errors.New("database is not connected")
)

func init() {
	connectToPostgres()
}

// connect to postgres and make automigrate
func connectToPostgres() {
	var err error
	dbObj, err = db.Connect(config)

	if err != nil {
		if errors.Is(err, db.ErrInvalidConfig) {
			log.Printf("[db]: Config is invalid: %q\n", config)
            return
		}
        // need to process other errors

        log.Printf("[db]: Got unknown error on connect: %q\n", err)
		return
	}

	log.Printf("[db]: Database connected\n")
}

