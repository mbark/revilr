package main

import (
	"code.google.com/p/go.crypto/bcrypt"
)

type User struct {
	Username string
	Password []byte
}


func (u *User) SetPassword(password string) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// you're in deep shit, son.
		panic(err)
	}
	u.Password = hashedPass;
}

func Login(username, password string) (user *User, err error) {
	user, err = getUser(username)

	if err != nil {
		user = nil
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		user = nil
	}
	return
}