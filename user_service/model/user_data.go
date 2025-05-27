package model

import "golang.org/x/crypto/bcrypt"


type User struct {
    ID       int64 `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// ValidCredentials - simply validate credentials for non-empty
func (u *User) ValidCredentials() bool {
    return u.Login != "" && u.Password != ""
}

func (u *User) HashPassword() error {
    bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)

    if err != nil {
        return err
    }
    
    u.Password = string(bytes)
    return nil
}

// ComparePassword compare u's plain password and given hashed password
func (u *User) ComparePassword(hashed string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(u.Password))
}
