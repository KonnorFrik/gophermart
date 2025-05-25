package model

import (
	"errors"
	"fmt"
	db "gophermart/model/postgres"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
    ErrConstraintUniqueViolation = errors.New("unique constraint violation")
    ErrConstraintForeignKeyViolation = errors.New("foreign key violation")
    ErrConstraintNotNullViolation = errors.New("not null violation")
    ErrConstraintCheckViolation = errors.New("check violation")
    ErrUnknown = errors.New("unknown error")
)


func init() {
	connectToPostgres()
}

// connect to postgres
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

func WrapError(err error) error {
    if err == nil {
        return nil
    } 

    var pgErr *pgconn.PgError

    if errors.As(err, &pgErr) {
        switch pgErr.Code {
        case "23505":
            return fmt.Errorf("%w: %q", ErrConstraintUniqueViolation, err)
        case "23503":
            return fmt.Errorf("%w: %q", ErrConstraintForeignKeyViolation, err)
        case "23502":
            return fmt.Errorf("%w: %q", ErrConstraintNotNullViolation, err)
        case "23514": 
            return fmt.Errorf("%w: %q", ErrConstraintCheckViolation, err)

        case "57P01":
            return fmt.Errorf("%w: %q", ErrDataBaseNotConnected, err)
        case "57P02":
            return fmt.Errorf("%w: %q", ErrDataBaseNotConnected, err)
        }
    }

    return fmt.Errorf("%w: %q", ErrUnknown, err)
}
