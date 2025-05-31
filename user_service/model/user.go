package model

import (
	"context"
	psgs "gophermart/db/postgres"
	"gophermart/internal/generated/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// TODO: write in doc-comment what data need to be already stored in struct for operation

type User struct {
    ID       int64 `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Create - create a new user in db from u's data.
// After successfull create fill 'u' with created data.
func (u *User) Create(ctx context.Context) error {
    if !u.validCredentials() {
        return ErrInvalidData
    }
    
    if err := u.hashPassword(); err != nil {
        return ErrInvalidData
    }

    db := psgs.DB()
    var err error
    userDB, err := db.CreateUser(ctx, models.CreateUserParams{
        Login: u.Login,
        Password: u.Password,
    })

	if err != nil {
		log.Printf("[model.User.NewUser]: Error on Create: %q\n", err)
		return wrapError(db.WrapError(err))
	}

    u.ID = userDB.ID
    u.Password = userDB.Password
    return nil
}

// ByCreadentials - retrieve user from db by credentials (login, password) from u.
// After successfull search fill 'u' with fetched data.
func (u *User) ByCreadentials(ctx context.Context) error {
    if !u.validCredentials() {
        return ErrInvalidData
    }

    db := psgs.DB()
    userDB, err := db.UserByLogin(ctx, u.Login)

	if err != nil {
		log.Printf("[model.User.UserByCredentials]: Error: %q\n", err)
		return wrapError(db.WrapError(err))
	}

    if err := u.comparePassword(userDB.Password); err != nil {
        log.Printf("[model.User/UserByCredentials]: Error on password compare: %q\n", err)
        return ErrInvalidData
    }

    u.ID = userDB.ID
    u.Password = userDB.Password
    return nil
}

// DeleteByID - delete a user from db by ID from 'u'.
func (u *User) DeleteByID(ctx context.Context) error {
    db := psgs.DB()
    err := db.DeleteUser(ctx, u.ID)

	if err != nil {
		log.Printf("[model.User/DeleteUser]: Error on delete %q\n", err)
		return wrapError(db.WrapError(err))
	}

    return nil
}


// validCredentials - simply validate credentials for non-empty
func (u *User) validCredentials() bool {
    return u.ID >= 0 && u.Login != "" && u.Password != ""
}

func (u *User) hashPassword() error {
    bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

    if err != nil {
        return err
    }
    
    u.Password = string(bytes)
    return nil
}

// ComparePassword compare u's plain password and given hashed password
func (u *User) comparePassword(hashed string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(u.Password))
}
