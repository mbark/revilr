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
	user := getUser(request)
	if user == nil {
		return
	}
	revils, err := db.GetRevilsOfType(revilType, *user)
	if err != nil {
		fmt.Println(err)
		revils = make([]data.Revil, 0)
	}

	revilsMap := RevilsAsMap(revils)
	ShowResponsePage(writer, request, revilType, revilsMap)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	if !isLoggedIn(request) {
		http.Redirect(writer, request, "/login", http.StatusMovedPermanently)
	} else {
		user := getUser(request)
		revils := make([]data.Revil, 0)
		if user != nil {
			allRevils, err := db.GetAllRevilsInDatabase(*user)
			if err != nil {
				fmt.Println(err)
			} else {
				revils = allRevils
			}
		}

		revilsMap := RevilsAsMap(revils)
		ShowResponsePage(writer, request, "home", revilsMap)		
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
		ShowResponsePage(writer, request, "login", make(map[string]interface{}))
	}
}

func userHandler(writer http.ResponseWriter, request *http.Request) {
	user := getUser(request)
	if user != nil {
		ShowResponsePage(writer, request, "user", user.AsMap())
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
		ShowResponsePage(writer, request, "register", make(map[string]interface{}))
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
		ShowResponsePage(writer, request, "logout", make(map[string]interface{}))
	}

}

func revilHandler(writer http.ResponseWriter, request *http.Request) {
	if !isLoggedIn(request) {
		http.Redirect(writer, request, "/login", http.StatusFound)
	} else {
		ShowResponsePage(writer, request, "revil", make(map[string]interface{}))
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

func getUser(request *http.Request) (user *data.User) {
	userId, ok := getUserId(request)
	if !ok {
		return
	}

	user, err := db.FindUserById(userId)
	if err != nil {
		fmt.Println(err)
		user = nil
	}
	return
}

func ShowResponsePage(writer http.ResponseWriter, request *http.Request, name string, data map[string]interface{}) {
	user := getUser(request)
	html := RenderWithAdditionalData(name, user, data)
	fmt.Fprintf(writer, html)
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