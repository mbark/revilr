package main

import (
	"encoding/gob"
	"net/http"
	"regexp"
	"revilr/db"
	"revilr/user"
)

const lenPath = len("/revilr/")

var validTypes = regexp.MustCompile("^(page|image|selection)$")

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
	DisplayRevils(revils, revilType, writer, request)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	revils := db.GetAllRevilsInDatabase()
	DisplayRevils(revils, "all", writer, request)
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	session := getSession(request)

	if session.Values["user"] != nil {
		http.Redirect(writer, request, "/user", http.StatusFound)
	}

	if request.Method == "POST" {
		username, password := parseUser(request)

		user := verifyUser(username)
		if user == nil {
			DisplayLogin(writer, "invalidUsername", request)
			return
		}

		loggedIn := user.Login(password)
		if loggedIn {
			session.Values["user"] = user
			err := session.Save(request, writer)
			if err == nil {
				http.Redirect(writer, request, "/user", http.StatusFound)
			} else {
				panic(err)
			}

		} else {
			DisplayLogin(writer, "invalidPassword", request)
		}
	} else if request.Method == "GET" {
		DisplayLogin(writer, "", request)
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
	loggedIn := getLoggedIn(request)

	if loggedIn {
		DisplayUser(writer, request)
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
			DisplayRegister(writer, "usernameTaken", request)
			return
		}
		user := &user.User{Username: username}
		user.SetPassword(password)
		err := db.CreateUser(user)

		if err != nil {
			DisplayRegister(writer, "failed", request)
			return
		}

		http.Redirect(writer, request, "/user", http.StatusFound)
		return
	} else if request.Method == "GET" {
		DisplayRegister(writer, "", request)
	}
}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {
	loggedIn := !logOut(writer, request)
	DisplayLogout(writer, loggedIn, request)
}
