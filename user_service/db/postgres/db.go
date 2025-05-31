package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"gophermart/internal/generated/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DbConn struct {
    conn *pgx.Conn
    *models.Queries
}

type DbConfig struct {
    Host string
    User string 
    Password string 
    DbName string 
    Port string
    SSLMode bool
    TimeZone string
}

var (
    // ErrInvalidConfig - invalid configuration for connect to db
    ErrInvalidConfig = errors.New("invalid database configuration")
    // ErrDataBaseNotConnected - no connection to db
    ErrDataBaseNotConnected = errors.New("database is not connected")
    // ErrConstraintUniqueViolation - dublication in db
    ErrConstraintUniqueViolation = errors.New("unique constraint violation")
    // ErrConstraintForeignKeyViolation - db error
    ErrConstraintForeignKeyViolation = errors.New("foreign key violation")
    // ErrConstraintNotNullViolation - null value in db
    ErrConstraintNotNullViolation = errors.New("not null violation")
    // ErrConstraintCheckViolation - custom check not passed
    ErrConstraintCheckViolation = errors.New("check violation")
    // ErrUnknown - any other not documented error in db
    ErrUnknown = errors.New("unknown error")

    // Default configuration for connect to postgres container from this repo
	DefaultConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DbName:   os.Getenv("DB_DB"),
		Port:     os.Getenv("DB_PORT"),
		// TimeZone: "Asia/Yekaterinburg",
	}

    dbObj DbConn
)

func init() {
    err := dbObj.connect()

    if err != nil {
        log.Printf("[db]: Connection error: %q\n", err)
        return 
    }

    log.Printf("[db]: DataBase connected\n")
}

// DB - return a single object for interact with postgres db
func DB() *DbConn {
    return &dbObj
}

// WrapError - wraps an error from the db into one of the predefined ones
func (dc *DbConn) WrapError(err error) error {
    newErr := wrapError(err)

    if newErr == nil {
        return nil
    }

    if errors.Is(newErr, ErrDataBaseNotConnected) {
        dc.connect()
    }

    return newErr
}

func (dc *DbConn) connect() error {
    log.Println("[db/.connect]: start connect to postgres")

    if !DefaultConfig.isValid() {
        log.Printf("Invalid config: %+v\n", DefaultConfig)
        return ErrInvalidConfig
    }

    var db *pgx.Conn
    var err error
    connectionInfo := DefaultConfig.String()
    db, err = pgx.Connect(context.Background(), connectionInfo)

    if err != nil {
        wrapped := wrapError(err)
        log.Printf("[db/.connect]: Error on connection: %q\n", wrapped)
        return wrapped
    }

    dc.conn = db
    dc.Queries = models.New(dc.conn)
    log.Println("[db/.connect]: Postgres Connected")
    return nil
}

func wrapError(err error) error {
    if err == nil {
        return nil
    } 

    var pgErr *pgconn.PgError

    if errors.As(err, &pgErr) {
        switch pgErr.Code {
        case "23505":
            return fmt.Errorf("%w: %w", ErrConstraintUniqueViolation, err)
        case "23503":
            return fmt.Errorf("%w: %w", ErrConstraintForeignKeyViolation, err)
        case "23502":
            return fmt.Errorf("%w: %w", ErrConstraintNotNullViolation, err)
        case "23514": 
            return fmt.Errorf("%w: %w", ErrConstraintCheckViolation, err)

        case "57P01":
            return fmt.Errorf("%w: %w", ErrDataBaseNotConnected, err)
        case "57P02":
            return fmt.Errorf("%w: %w", ErrDataBaseNotConnected, err)

        default:
            log.Printf("[db] Can't process code: %q\n", pgErr.Code)
        }
    }

    return fmt.Errorf("%w: %w", ErrUnknown, err)
}

func (dc *DbConfig) String() string {
    sslMode := "disable"

    if dc.SSLMode {
        sslMode = "enable"
    }

    return fmt.Sprintf(
        "host=%s database=%s user=%s password=%s port=%s sslmode=%s",
        dc.Host,
        dc.DbName,
        dc.User,
        dc.Password,
        dc.Port,
        sslMode,
        // dc.TimeZone,
    )
}

func (dc *DbConfig) isValid() bool {
    return dc.Host != "" && dc.User != "" && dc.Password != "" && dc.Port != ""// && dc.TimeZone != "" 
}
