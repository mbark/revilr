package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"revilr/data"
)

var store = sessions.NewCookieStore([]byte("this-is-a-secret"))
var user_session = "users"

func getSession(request *http.Request) *sessions.Session {
	session, err := store.Get(request, user_session)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return session
}

func getUser(request *http.Request) (bool, data.User) {
	user, ok := getSession(request).Values["user"].(data.User)
	return ok, user
}

func logOut(writer http.ResponseWriter, request *http.Request) bool {
	session := getSession(request)

	session.Values["user"] = nil
	session.Save(request, writer)
	if session.Values["user"] == nil {
		return true
	}
	return false
}
