package model

import "errors"

var (
    users = map[string]*User{}
    userToken = map[string]string{} // login to token
)

var (
    ErrUserAlreadyExist = errors.New("user already exist")
    ErrUserDoesNotExist = errors.New("user does not exist")
)

type User struct {
    Login string
    Password string
    Token string
}

func NewUser(login, passw, tok string) (*User, error) {
    if _, exist := users[login]; exist {
        return nil, ErrUserAlreadyExist
    }

    user := &User{
        Login: login,
        Password: passw,
        Token: tok,
    }

    users[login] = user
    return user, nil
}

func UserByLogin(login string) (*User, error) {
    user, ok := users[login]

    if !ok {
        return nil, ErrUserDoesNotExist
    }

    return user, nil
}

func DeleteUser(login string) {
    delete(users, login)
}
