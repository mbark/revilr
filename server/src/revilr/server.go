package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"revilr/data"
	"revilr/db"
)

const lenPath = len("/revilr/")

var validTypes = regexp.MustCompile("^(page|image|selection)$")
var validEmail = regexp.MustCompile("([a-z]*.)+@[a-z]+.[a-z]+")

func main() {
	db.OpenConnection()

	http.HandleFunc("/revilr/", httpHandler)
	http.HandleFunc("/revilr", indexHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/revil", revilHandler)
	http.HandleFunc("/user_taken", userTakenHandler)
	http.HandleFunc("/email_taken", emailTakenHandler)
	http.HandleFunc("/user_valid", isValidUserHandler)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	http.ListenAndServe(":8080", nil)
}

func httpHandler(writer http.ResponseWriter, request *http.Request) {
	loggedIn := isLoggedIn(request)
	if !loggedIn {
		http.Redirect(writer, request, "/login", http.StatusMovedPermanently)
	}

	revilType, success := parseType(request)
	if !success {
		http.NotFound(writer, request)
		return
	}

	if request.Method == "POST" {
		postHandler(request, revilType)
		http.Redirect(writer, request, "/revilr", http.StatusTemporaryRedirect)
	} else if request.Method == "GET" {
		getHandler(writer, request, revilType)
	}
}

func parseType(request *http.Request) (string, bool) {
	revilType := request.URL.Path[lenPath:]
	if !validTypes.MatchString(revilType) {
		revilType = request.FormValue("type")
	}
	return revilType, validTypes.MatchString(revilType)
}

func postHandler(request *http.Request, revilType string) {
	userId, ok := getUserId(request)
	if !ok {
		return
	}

	url := request.FormValue("url")
	title := request.FormValue("title")
	note := request.FormValue("note")

	err := db.CreateRevil(userId, revilType, url, title, note)
	if err != nil {
		fmt.Println(err)
	}
}

func getHandler(writer http.ResponseWriter, request *http.Request, revilType string) {
	user, ok := getUser(request)
	if !ok {
		return
	}
	revils, err := db.GetRevilsOfType(revilType, user)
	if err != nil {
		fmt.Println(err)
		revils = make([]data.Revil, 0)
	}
	DisplayRevils(revils, revilType, writer, request)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	if !isLoggedIn(request) {
		http.Redirect(writer, request, "/login", http.StatusMovedPermanently)
	} else {
		user, ok := getUser(request)
		revils := make([]data.Revil, 0)
		if ok {
			allRevils, err := db.GetAllRevilsInDatabase(user)
			if err != nil {
				fmt.Println(err)
			} else {
				revils = allRevils
			}
		}

		DisplayRevils(revils, "home", writer, request)
	}
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	if isLoggedIn(request) {
		http.Redirect(writer, request, "/user", http.StatusTemporaryRedirect)
	}

	if request.Method == "POST" {
		user, canLogin := verifyUser(request)

		if canLogin && user != nil {
			if setUser(writer, request, *user) == nil {
				http.Redirect(writer, request, "/user", http.StatusMovedPermanently)
			}
		}
	} else if request.Method == "GET" {
		DisplayLogin(writer, request)
	}
}

func userHandler(writer http.ResponseWriter, request *http.Request) {
	user, ok := getUser(request)

	if ok {
		DisplayUser(writer, request, user)
	} else {
		http.Redirect(writer, request, "/login", http.StatusMovedPermanently)
	}
}

func registerHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		if isValidRegister(request) {
			username := request.FormValue("username")
			password := request.FormValue("password")
			email := request.FormValue("email")

			user, err := db.CreateUser(username, password, email)
			if err == nil {
				err = setUser(writer, request, user)
				if err == nil {
					http.Redirect(writer, request, "/user", http.StatusTemporaryRedirect)
					return
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Invalid register")
		}
	} else if request.Method == "GET" {
		DisplayRegister(writer, request)
	}
}

func isValidRegister(request *http.Request) bool {
	username := request.FormValue("username")
	password := request.FormValue("password")
	verification := request.FormValue("verification")
	email := request.FormValue("email")

	if len(username) < 5 || len(username) > 12 {
		return false
	}
	if !validEmail.MatchString(email) {
		return false;
	}
	if len(password) < 8 {
		return false
	}
	if password != verification {
		return false
	}
	if tmp, _ := verifyUser(request); tmp != nil {
		return false
	}
	if tmp, _ := db.FindUserByEmail(email); tmp != nil {
		return false;
	}

	return true
}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {
	err := logOut(writer, request)
	if err != nil {
		fmt.Println(err)
	} else {
		DisplayLogout(writer, request)
	}

}

func revilHandler(writer http.ResponseWriter, request *http.Request) {
	if !isLoggedIn(request) {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		DisplayRevil(writer, request)
	}
}

func emailTakenHandler(writer http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")
	user, _ := db.FindUserByEmail(email)
	isTaken := user != nil

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, Response{"isTaken": isTaken})
}

func userTakenHandler(writer http.ResponseWriter, request *http.Request) {
	user, _ := verifyUser(request)
	isTaken := user != nil

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, Response{"isTaken": isTaken})
}

func isValidUserHandler(writer http.ResponseWriter, request *http.Request) {
	isValid := false
	_, canLogin := verifyUser(request)

	if canLogin {
		isValid = true
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, Response{"isValid": isValid})
}

func verifyUser(request *http.Request) (*data.User, bool) {
	username := request.FormValue("username")
	password := request.FormValue("password")

	user, err := db.FindUserByName(username)
	if err == nil && user != nil {
		canLogin := user.PasswordMatches(password)
		return user, canLogin
	}
	return nil, false
}

func getUser(request *http.Request) (user data.User, ok bool) {
	userId, ok := getUserId(request)
	if !ok {
		return
	}

	userPointer, err := db.FindUserById(userId)
	ok = err == nil
	if userPointer != nil {
		user = *userPointer
	}
	return
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
