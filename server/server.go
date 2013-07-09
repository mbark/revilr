package main

import (
	"fmt"
	"net/http"
	"regexp"
)

const lenPath = len("/revilr/")

var validTypes = regexp.MustCompile("^(page|image|selection)$")

func main() {
	db, err := getDatabase()
	if err != nil {
		return
	}

	database = db
	defer database.Close()

	http.HandleFunc("/revilr/", httpHandler)
	http.HandleFunc("/revilr", indexHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
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
	if request.Method == "POST" {
		username, password := parseUser(request)

		user := verifyUser(username)
		if user == nil {
			DisplayLogin(writer, "invalidUsername")
			return
		}

		loggedIn := user.Login(password)
		if loggedIn {
			http.Redirect(writer, request, "/user", http.StatusFound)
		} else {
			DisplayLogin(writer, "invalidPassword")
			return
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
	DisplayUser(writer, "hej")
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
			fmt.Println("username is taken")
			DisplayRegister(writer, "usernameTaken")
			return
		}
		user := &User{Username: username}
		user.SetPassword(password)
		err := createUser(user)

		if err != nil {
			fmt.Println("failed to create...")
			DisplayRegister(writer, "failed")
			return
		}

		fmt.Println("successfull!")
		http.Redirect(writer, request, "/user", http.StatusFound)
		return
	} else if request.Method == "GET" {
		DisplayRegister(writer, "")
	}
}
