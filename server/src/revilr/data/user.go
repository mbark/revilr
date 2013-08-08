package data

import (
	"code.google.com/p/go.crypto/bcrypt"
	"labix.org/v2/mgo/bson"
)

func CreateUser(username string, password, email string) (user User, err error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user = User {
		Id: bson.NewObjectId(),
		Username: username,
		Email: email,
		Password: hashedPass,
		Created: bson.Now(),
	}
	return
}

func (user User) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return (err == nil)
}
