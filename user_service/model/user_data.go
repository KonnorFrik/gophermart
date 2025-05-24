package model


type User struct {
    ID       int64 `json:"-"`
    Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// ValidCredentials - simply validate credentials for non-empty
func (u *User) ValidCredentials() bool {
    return u.Email != "" && u.Login != "" && u.Password != ""
}
