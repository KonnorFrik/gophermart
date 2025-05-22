package model

import (
	"errors"
	db "gophermart/model/postgres"
	"log"
	"os"

    // pgx
    // sqlc - https://github.com/sqlc-dev/sqlc
    // https://github.com/doug-martin/goqu
    // migrations
	"gorm.io/gorm"
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

	dbObj *gorm.DB
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
			log.Printf("[model]: Config is invalid: %q\n", config)

		} else {
			// need to ...
		}

		return
	}

	if err = dbObj.AutoMigrate(&User{}); err != nil {
		log.Printf("[model]: Error on migrate User: %q\n", err)
		return
	}
	
	if err = dbObj.AutoMigrate(&Order{}); err != nil {
		log.Printf("[model]: Error on migrate Order: %q\n", err)
		return
	}

	log.Printf("[model]: Database connected\n")
}

