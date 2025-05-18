package model

import (
	"errors"
	"log"
	"os"

	db "gophermart/model/postgres"

	"gorm.io/gorm"
)

var (
	config = db.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DbName:   os.Getenv("DB_DB"),
		Port:     os.Getenv("DB_PORT"),
		TimeZone: "Asia/Yekaterinburg",
	}

	dbObj *gorm.DB
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserDoesNotExist = errors.New("user does not exist")
    ErrDataBaseNotConnected = errors.New("database is not connected")
)

type User struct {
	Login    string `gorm:"unique"`
	Password string
	Token    string
}

func init() {
	connectToPostgres()
}

// connect to postgres and make automigrate
func connectToPostgres() {
	var err error
	dbObj, err = db.Connect(config)

	if err != nil {
		if errors.Is(err, db.ErrInvalidConfig) {
			log.Printf("[model:User]: Config is invalid: %q\n", config)

		} else {
			// need to run reconnector
		}
		return
	}

	err = dbObj.AutoMigrate(User{})

	if err != nil {
		log.Printf("[model:User]: Error on migrate: %q\n", err)
		return
	}

	log.Printf("[model:User]: Database connected\n")
}

func NewUser(login, passw, tok string) (*User, error) {
    if dbObj == nil {
        connectToPostgres()
        return nil, ErrDataBaseNotConnected
    }

    var user User = User{
        Login: login,
        Password: passw,
        Token: tok,
    }
    tx := dbObj.Model(User{}).FirstOrCreate(&user, user)

    if tx.Error != nil {
        log.Printf("[model:User/NewUser]: Error on FirstOrCreate: %q\n", tx.Error)
        return nil, tx.Error
    }

    return &user, nil
}

func UserByLogin(login string) (*User, error) {
    if dbObj == nil {
        connectToPostgres()
        return nil, ErrDataBaseNotConnected
    }

    var user User
    tx := dbObj.Model(User{}).First(&user, "login = ?", login)

    if tx.Error != nil {
        log.Printf("[model:User/UserByLogin]: Error: %q\n", tx.Error)

        if tx.Error == gorm.ErrRecordNotFound {
            return nil, ErrUserDoesNotExist
        }

        return nil, tx.Error
    }

    return &user, nil
}

func DeleteUser(user *User) error {
    if dbObj == nil {
        connectToPostgres()
        return ErrDataBaseNotConnected
    }

    tx := dbObj.Model(User{}).Delete(user, "login = ?", user.Login)

    if tx.Error != nil {
        log.Printf("[model:User/DeleteUser]: Error on delete %q\n", tx.Error)
        return tx.Error
    }

    return nil
}
