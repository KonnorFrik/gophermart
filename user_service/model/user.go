package model

import (
	"context"
	"errors"
	"gophermart/model/models"
	"log"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserDoesNotExist = errors.New("user does not exist")
)

func NewUser(login, passw string) (*models.User, error) {
	if dbObj == nil {
		log.Printf("[model.User/NewUser]: Lost connection to DB\n")
		connectToPostgres()
		return nil, ErrDataBaseNotConnected
	}

	var err error
    queries := models.New(dbObj)

    user, err := queries.CreateUser(context.Background(), models.CreateUserParams{
        Login: login,
        Password: passw,
    })

	if err != nil {
		log.Printf("[model.User/NewUser]: Error on Create: %q\n", err)
		return nil, err
	}

	return user, nil
}

func UserByLogin(login string) (*models.User, error) {
	if dbObj == nil {
		log.Printf("[model.User/UserByLogin]: Lost connection to DB\n")
		connectToPostgres()
		return nil, ErrDataBaseNotConnected
	}

    queries := models.New(dbObj)
    user, err := queries.UserByLogin(context.Background(), login)

	if err != nil {
		log.Printf("[model.User/UserByLogin]: Error: %q\n", err)
		return nil, err
	}

	return user, nil
}

func DeleteUser(id int64) error {
	if dbObj == nil {
		log.Printf("[model.User/DeleteUser]: Lost connection to DB\n")
		connectToPostgres()
		return ErrDataBaseNotConnected
	}

    queries := models.New(dbObj)
    err := queries.DeleteUser(context.Background(), id)

	if err != nil {
		log.Printf("[model.User/DeleteUser]: Error on delete %q\n", err)
		return err
	}

	return nil
}
