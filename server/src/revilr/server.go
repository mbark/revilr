package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
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
	http.HandleFunc("/revil", revilHandler)
	http.HandleFunc("/user_taken", userTakenHandler)
	http.HandleFunc("/user_valid", isValidUserHandler)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe(":8080", nil)
}

func httpHandler(writer http.ResponseWriter, request *http.Request) {
	loggedIn, usr := getUser(request)
	if !loggedIn {
		http.Redirect(writer, request, "/login", http.StatusMovedPermanently)
	}

	revilType, success := parseType(request)
	if !success {
		http.NotFound(writer, request)
		return
	}

	if request.Method == "POST" {
		postHandler(request, revilType, usr)
		http.Redirect(writer, request, "/revilr", http.StatusTemporaryRedirect)
	} else if request.Method == "GET" {
		getHandler(writer, request, revilType, usr)
	}
}

func parseType(request *http.Request) (string, bool) {
	revilType := request.URL.Path[lenPath:]
	if !validTypes.MatchString(revilType) {
		revilType = request.FormValue("type")
	}
	return revilType, validTypes.MatchString(revilType)
}

func postHandler(request *http.Request, revilType string, usr user.User) {
	rev := user.Revil{Type: revilType, Url: request.FormValue("url"), Comment: request.FormValue("c")}
	rev.PrintRevil()
	err := db.InsertIntoDatabase(rev, usr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revild ", rev)
	}

}

func getHandler(writer http.ResponseWriter, request *http.Request, revilType string, usr user.User) {
	revils := db.GetRevilsOfType(revilType, usr)
	DisplayRevils(revils, revilType, writer, request)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	loggedIn, user := getUser(request)
	if !loggedIn {
		http.Redirect(writer, request, "/login", http.StatusMovedPermanently)
	} else {
		revils := db.GetAllRevilsInDatabase(user)
		DisplayRevils(revils, "all", writer, request)
	}
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	session := getSession(request)

	if session.Values["user"] != nil {
		http.Redirect(writer, request, "/user", http.StatusTemporaryRedirect)
	}

	if request.Method == "POST" {
		username, password := parseUser(request)

		user := verifyUser(username)
		if user != nil {

			loggedIn := user.Login(password)
			if loggedIn {
				session.Values["user"] = user
				err := session.Save(request, writer)
				if err == nil {
					http.Redirect(writer, request, "/user", http.StatusMovedPermanently)
				} else {
					panic(err)
				}

			}
		}
	} else if request.Method == "GET" {
		DisplayLogin(writer, request)
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
	loggedIn, _ := getUser(request)

	if loggedIn {
		DisplayUser(writer, request)
	} else {
		http.Redirect(writer, request, "/login", http.StatusMovedPermanently)
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
		if verifyUser(username) == nil {
			user := &user.User{Username: username}
			user.SetPassword(password)
			err := db.CreateUser(user)

			if err == nil {
				session := getSession(request)
				session.Values["user"] = user
				err := session.Save(request, writer)
				if err == nil {
					http.Redirect(writer, request, "/user", http.StatusTemporaryRedirect)
					return
				}
			}
		}
	} else if request.Method == "GET" {
		DisplayRegister(writer, request)
	}
}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {
	loggedIn := !logOut(writer, request)
	DisplayLogout(writer, loggedIn, request)
}

func revilHandler(writer http.ResponseWriter, request *http.Request) {
	loggedIn, _ := getUser(request)
	if !loggedIn {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		DisplayRevil(writer, request)
	}
}

func userTakenHandler(writer http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	user := verifyUser(username)
	isTaken := user != nil

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, Response{"isTaken": isTaken})
}

func isValidUserHandler(writer http.ResponseWriter, request *http.Request) {
	isValid := false

	username, password := parseUser(request)
	user := verifyUser(username)

	if user != nil {
		canLogin := user.Login(password)
		if canLogin {
			isValid = true
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, Response{"isValid": isValid})
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}
