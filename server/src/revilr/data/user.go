package data

import (
	"code.google.com/p/go.crypto/bcrypt"
	"labix.org/v2/mgo/bson"
)

func CreateUser(username string, password string) (user User, err error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user = User{Username: username, Password: hashedPass, Id: bson.NewObjectId(), Created: bson.Now()}
	return
}

func (user User) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return (err == nil)
}
