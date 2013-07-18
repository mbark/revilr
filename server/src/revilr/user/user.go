package user

import (
	"code.google.com/p/go.crypto/bcrypt"
)

type User struct {
	Username string
	Password []byte
}

func NewType(username string, password []byte) *User {
	return &User{Username: username, Password: password}
}

func (u *User) SetPassword(password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = hashedPass
	return nil
}

func (user *User) Login(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return (err == nil)
}