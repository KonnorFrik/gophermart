package model


type User struct {
    ID       int64 `json:"-"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// ValidCredentials - simply validate credentials for non-empty
func (u *User) ValidCredentials() bool {
    return u.Login != "" && u.Password != ""
}
