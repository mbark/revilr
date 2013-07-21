package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"revilr/user"
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

func getLoggedIn(request *http.Request) bool {
	return getSession(request).Values["user"] != nil
}

func getUsername(request *http.Request) (bool, string) {
	user, ok := getSession(request).Values["user"].(user.User)
	if ok {
		return true, user.Username
	}
	return false, ""
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
