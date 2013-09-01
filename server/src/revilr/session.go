package main

import (
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
	"revilr/data"
)

const userSession = "users"

var (
	ErrUserIdOfWrongType = errors.New("User id is not of expected type")
	ErrUserNogLoggedIn = errors.New("User is not logged in")
)

var store = sessions.NewCookieStore([]byte(""))

func getSession(request *http.Request) (session *sessions.Session, err error) {
	session, err = store.Get(request, userSession)
	return
}

func setUser(writer http.ResponseWriter, request *http.Request, user data.User) (err error) {
	session, err := getSession(request)
	if err != nil {
		return
	}

	session.Values["uid"] = user.Id.Hex()
	err = session.Save(request, writer)
	return
}

func getUserId(request *http.Request) (userId string, err error) {
	session, err := getSession(request)
	if err != nil {
		return
	}

	tmp := session.Values["uid"]
	if tmp == nil {
		err = ErrUserNogLoggedIn
		return
	}

	userId, ok := tmp.(string)

	if !ok {
		err = ErrUserIdOfWrongType
	}

	return
}

func logOut(writer http.ResponseWriter, request *http.Request) (err error) {
	session, err := getSession(request)
	if err != nil {
		return
	}

	session.Values["uid"] = nil
	err = session.Save(request, writer)
	return err
}
