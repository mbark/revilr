package main

import (
	"github.com/gorilla/sessions"
	"net/http"
	"revilr/data"
)

const userSession = "users"

var store = sessions.NewCookieStore([]byte(""))

func getSession(request *http.Request) *sessions.Session {
	session, err := store.Get(request, userSession)
	if err != nil {
		panic(err)
	}
	return session
}

func setUser(writer http.ResponseWriter, request *http.Request, user data.User) error {
	session := getSession(request)
	session.Values["uid"] = user.Id.Hex()
	err := session.Save(request, writer)
	return err
}

func getUserId(request *http.Request) (userId string, ok bool) {
	userId, ok = getSession(request).Values["uid"].(string)
	return
}

func isLoggedIn(request *http.Request) bool {
	return (getSession(request).Values["uid"] != nil)
}

func logOut(writer http.ResponseWriter, request *http.Request) error {
	session := getSession(request)
	session.Values["uid"] = nil
	err := session.Save(request, writer)
	return err
}
