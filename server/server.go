package main

import (
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
	db, err := getDatabase()
	if err != nil {
		panic(err)
	}

	database = db
	defer database.Close()

	gob.Register(User{})

	http.HandleFunc("/revilr/", httpHandler)
	http.HandleFunc("/revilr", indexHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("templates/resources"))))
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
	rev := revil{Type: revilType, Url: request.FormValue("url"), Comment: request.FormValue("c")}
	rev.printRevil()
	insertIntoDatabase(rev)
}

func getHandler(writer http.ResponseWriter, request *http.Request, revilType string) {
	revils := getRevilsOfType(revilType)
	DisplayRevils(revils, revilType, writer)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	revils := getAllRevilsInDatabase()
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

func verifyUser(username string) *User {
	if username != "" {
		user, err := findUser(username)
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
	user, ok := session.Values["user"].(User)
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
		user := &User{Username: username}
		user.SetPassword(password)
		err := createUser(user)

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
