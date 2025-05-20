package model

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserDoesNotExist = errors.New("user does not exist")
)

type User struct {
    ID       uint    `gorm:"primarykey"`
	Login    string  `gorm:"unique"`
	Password string
	Orders   []Order `gorm:"constraint:OnDelete:CASCADE;"`
}

// TODO: pass a User type as arg
// and return just a error
func NewUser(login, passw string) (*User, error) {
	if dbObj == nil {
		log.Printf("[model.User/NewUser]: Lost connection to DB\n")
		connectToPostgres()
		return nil, ErrDataBaseNotConnected
	}

	var err error

	// var dummy User
	// err = dbObj.Model(&User{}).First(&dummy).Error
	//
	// if err == nil {
	//     log.Printf("[model.User/NewUser]: User already exist: %q\n", dummy.Login)
	//     return nil, ErrUserAlreadyExist
	// }
	//
	// if !errors.Is(err, gorm.ErrRecordNotFound) {
	//     log.Printf("[model.User/NewUser]: Error on find user: %q\n", err)
	//     return nil, err
	// }

	var user User = User{
		Login:    login,
		Password: passw,
	}
	err = dbObj.Model(&User{}).Create(&user).Error

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Printf("[model.User/NewUser]: User already exist: %q\n", user.Login)
		return nil, ErrUserAlreadyExist
	}

	if err != nil {
		log.Printf("[model.User/NewUser]: Error on Create: %q\n", err)
		return nil, err
	}

	return &user, nil
}

func UserByLogin(login string) (*User, error) {
	if dbObj == nil {
		log.Printf("[model.User/UserByLogin]: Lost connection to DB\n")
		connectToPostgres()
		return nil, ErrDataBaseNotConnected
	}

	var user User
	var err error
	err = dbObj.Model(&User{}).First(&user, "login = ?", login).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[model.User/UserByLogin]: User does not exist: %q\n", login)
		return nil, ErrUserDoesNotExist
	}

	if err != nil {
		log.Printf("[model.User/UserByLogin]: Error: %q\n", err)
		return nil, err
	}

	return &user, nil
}

func DeleteUser(user *User) error {
	if dbObj == nil {
		log.Printf("[model.User/DeleteUser]: Lost connection to DB\n")
		connectToPostgres()
		return ErrDataBaseNotConnected
	}

	var err error
	err = dbObj.Model(&User{}).Delete(user, "login = ?", user.Login).Error

	if err != nil {
		log.Printf("[model.User/DeleteUser]: Error on delete %q\n", err)
		return err
	}

	return nil
}
