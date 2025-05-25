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
    ErrUserInvalidData = errors.New("invalid data for create user")
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
    
    // hash password here

    queries := getQueries()
    defer putQueries(queries)
    var err error

    userDB, err := queries.CreateUser(context.Background(), models.CreateUserParams{
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
        Password: userDB.Password,
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
    userDB, err := queries.UserByLogin(context.Background(), user.Login)
    err = WrapError(err)

    switch {
    case errors.Is(err, ErrDataBaseNotConnected):
        connectToPostgres()
        return nil, err
    }

	if err != nil {
		log.Printf("[model.User/UserByLogin]: Error: %q\n", err)
		return nil, err
	}

    // check hashed password somewhere here

    toRet := &User{
        ID: userDB.ID,
        Login: userDB.Login,
        Password: userDB.Password,
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
    err := queries.DeleteUser(context.Background(), id)
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
