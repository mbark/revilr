package main

import (
	"bytes"
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/gob"
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

func (user *User) GobEncode() ([]byte, error) {
	writer := new(bytes.Buffer)
	encoder := gob.NewEncoder(writer)
	err := encoder.Encode(user.Username)
	if err != nil {
		return nil, err
	}
	err = encoder.Encode(user.Password)
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func (user *User) GobDecode(buf []byte) error {
	reader := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&user.Username)
	if err != nil {
		return err
	}
	return decoder.Decode(&user.Password)
}
