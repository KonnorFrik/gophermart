package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const (
)

var (
    ErrInvalidConfig = errors.New("invalid database configuration")
)

// Connect - connect to db with specified config which must be not nil
func Connect(config *DbConfig) (*pgx.Conn, error) {
    if !config.isValid() {
        return nil, ErrInvalidConfig
    }

    var db *pgx.Conn
    var err error
    connectionInfo := config.String()
    db, err = pgx.Connect(context.Background(), connectionInfo)

    if err != nil {
        return nil, err
    }

    return db, nil
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

func (dc *DbConfig) String() string {
    sslMode := "disable"

    if dc.SSLMode {
        sslMode = "enable"
    }

    return fmt.Sprintf(
        "host=%s database=%s user=%s password=%s port=%s sslmode=%s TimeZone=%s",
        dc.Host,
        dc.DbName,
        dc.User,
        dc.Password,
        dc.Port,
        sslMode,
        dc.TimeZone,
    )
}

func (dc *DbConfig) isValid() bool {
    return dc.Host != "" && dc.User != "" && dc.Password != "" && dc.Port != "" && dc.TimeZone != "" 
}
