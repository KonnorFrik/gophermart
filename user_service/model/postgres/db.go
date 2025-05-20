package postgres

import (
    "errors"
    "fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
)

var (
    ErrInvalidConfig = errors.New("invalid database configuration")
)

// Connect - connect to db with specified config which must be not nil
func Connect(config *DbConfig) (*gorm.DB, error) {
    if !config.isValid() {
        return nil, ErrInvalidConfig
    }

    var db *gorm.DB
    var err error
    connectionInfo := config.String()
    db, err = gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})

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
