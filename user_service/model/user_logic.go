package model

import (
	"context"
	"errors"
	"gophermart/internal/generated/models"
	"log"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserDoesNotExist = errors.New("user does not exist")
    ErrUserInvalidData = errors.New("invalid data")
)

// NewUser create a new user in DB.
// Accept raw data from request body.
func NewUser(user User) (*User, error) {
	if dbObj == nil {
		connectToPostgres()
		return nil, ErrDataBaseNotConnected
	}

    if !user.ValidCredentials() {
        return nil, ErrUserInvalidData
    }
    
    if err := user.HashPassword(); err != nil {
        return nil, ErrUserInvalidData
    }

    queries := getQueries()
    defer putQueries(queries)
    var err error

    userDB, err := queries.CreateUser(context.TODO(), models.CreateUserParams{
        Login: user.Login,
        Password: user.Password,
    })
    err = WrapError(err)

    switch {
    case errors.Is(err, ErrDataBaseNotConnected):
        connectToPostgres()
        return nil, err
    }

	if err != nil {
		log.Printf("[model.User/NewUser]: Error on Create: %q\n", err)
		return nil, err
	}

    toRet := &User{
        ID: userDB.ID,
        Login: userDB.Login,
        // Password: userDB.Password,
    }
	return toRet, nil
}

// UserByCredentials - Returns a user data stored in DB
func UserByCredentials(user User) (*User, error) {
	if dbObj == nil {
		connectToPostgres()
		return nil, ErrDataBaseNotConnected
	}

    if !user.ValidCredentials() {
        return nil, ErrUserInvalidData
    }

    queries := getQueries()
    defer putQueries(queries)
    userDB, err := queries.UserByLogin(context.TODO(), user.Login)
    err = WrapError(err)

    switch {
    case errors.Is(err, ErrDataBaseNotConnected):
        connectToPostgres()
        return nil, err
    }

	if err != nil {
		log.Printf("[model.User/UserByCredentials]: Error: %q\n", err)
		return nil, err
	}

    if err := user.ComparePassword(userDB.Password); err != nil {
        log.Printf("[model.User/UserByCredentials]: Error on password compare: %q\n", err)
        return nil, ErrUserInvalidData
    }

    toRet := &User{
        ID: userDB.ID,
        Login: userDB.Login,
        // Password: userDB.Password,
    }
	return toRet, nil
}

// DeleteUserById - Delete user from DB
func DeleteUserById(id int64) error {
	if dbObj == nil {
		connectToPostgres()
		return ErrDataBaseNotConnected
	}

    queries := getQueries()
    defer putQueries(queries)
    err := queries.DeleteUser(context.TODO(), id)
    err = WrapError(err)

    switch {
    case errors.Is(err, ErrDataBaseNotConnected):
        connectToPostgres()
        return err
    }

	if err != nil {
		log.Printf("[model.User/DeleteUser]: Error on delete %q\n", err)
		return err
	}

	return nil
}
