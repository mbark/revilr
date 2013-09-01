package data

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"fmt"
	"io"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"strings"
)

func CreateUser(username string, password, email string) (user User, err error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	m := md5.New()
	io.WriteString(m, fmt.Sprintf("%s", rand.Intn(100)))
	verification := fmt.Sprintf("%x", m.Sum(nil))

	user = User{
		Id:           bson.NewObjectId(),
		Username:     username,
		Email:        email,
		Password:     hashedPass,
		Verification: verification,
		Verified:     false,
		Created:      bson.Now(),
	}
	return
}

func (user User) AsMap() map[string]interface{} {
	data := make(map[string]interface{})

	data["username"] = user.Username
	data["email"] = user.Email
	data["emailHash"] = user.EmailHash()

	return data
}

func (user User) EmailHash() string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(user.Email))
	return fmt.Sprintf("%x", m.Sum(nil))
}

func (user User) PasswordMatches(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	return (err == nil)
}
