package main

import (
	"revilr/db"
	"revilr/user"
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"regexp"
)

const lenPath = len("/revilr/")

var validTypes = regexp.MustCompile("^(page|image|selection)$")

var store = sessions.NewCookieStore([]byte("this-is-a-secret"))
var user_session = "users"

func main() {
	db, err := db.OpenConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	gob.Register(user.User{})

	http.HandleFunc("/revilr/", httpHandler)
	http.HandleFunc("/revilr", indexHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe(":8080", nil)
}

func httpHandler(writer http.ResponseWriter, request *http.Request) {
	revilType := parseType(request)
	if !validTypes.MatchString(revilType) {
		http.NotFound(writer, request)
		return
	}

	if request.Method == "POST" {
		postHandler(request, revilType)
	} else if request.Method == "GET" {
		getHandler(writer, request, revilType)
	}
}

func parseType(request *http.Request) string {
	return request.URL.Path[lenPath:]
}

func postHandler(request *http.Request, revilType string) {
	rev := user.Revil{Type: revilType, Url: request.FormValue("url"), Comment: request.FormValue("c")}
	rev.PrintRevil()
	db.InsertIntoDatabase(rev)
}

func getHandler(writer http.ResponseWriter, request *http.Request, revilType string) {
	revils := db.GetRevilsOfType(revilType)
	DisplayRevils(revils, revilType, writer)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	revils := db.GetAllRevilsInDatabase()
	DisplayRevils(revils, "all", writer)
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := store.Get(request, user_session)
	if err != nil {
		panic(err)
	}
	if session.Values["user"] != nil {
		http.Redirect(writer, request, "/user", http.StatusFound)
	}
	if request.Method == "POST" {
		username, password := parseUser(request)

		user := verifyUser(username)
		if user == nil {
			DisplayLogin(writer, "invalidUsername")
			return
		}

		loggedIn := user.Login(password)
		if loggedIn {
			session.Values["user"] = user
			err = session.Save(request, writer)
			if err == nil {
				http.Redirect(writer, request, "/user", http.StatusFound)
			} else {
				panic(err)
			}

		} else {
			DisplayLogin(writer, "invalidPassword")
		}
	} else if request.Method == "GET" {
		DisplayLogin(writer, "")
	}
}

func verifyUser(username string) *user.User {
	if username != "" {
		user, err := db.FindUser(username)
		if err == nil {
			if user.Username != "" {
				return user
			}
		}
	}
	return nil
}

func userHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := store.Get(request, user_session)
	if err != nil {
		//when will this happen?
		fmt.Println(err)
		return
	}
	user, ok := session.Values["user"].(user.User)
	if ok {
		DisplayUser(writer, user.Username)
	} else {
		http.Redirect(writer, request, "/login", http.StatusFound)
	}

}

func parseUser(request *http.Request) (username, password string) {
	username = request.FormValue("username")
	password = request.FormValue("password")
	return
}

func registerHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		username, password := parseUser(request)
		if verifyUser(username) != nil {
			DisplayRegister(writer, "usernameTaken")
			return
		}
		user := &user.User{Username: username}
		user.SetPassword(password)
		err := db.CreateUser(user)

		if err != nil {
			DisplayRegister(writer, "failed")
			return
		}

		http.Redirect(writer, request, "/user", http.StatusFound)
		return
	} else if request.Method == "GET" {
		DisplayRegister(writer, "")
	}
}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {
	session, err := store.Get(request, user_session)
	if err != nil {
		panic(err)
	}
	if request.Method == "POST" {
		session.Values["user"] = nil
		session.Save(request, writer)
	}
	var isLoggedOut string
	if session.Values["user"] == nil {
		isLoggedOut = "loggedOut"
	} else {
		isLoggedOut = ""
	}
	DisplayLogout(writer, isLoggedOut)
}
