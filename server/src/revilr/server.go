package main

import (
"encoding/json"
"fmt"
"github.com/gorilla/mux"
"net/http"
"regexp"
"revilr/data"
"revilr/db"
)

const lenPath = len("/revilr/")

var validEmail = regexp.MustCompile("([a-z]*.)+@[a-z]+.[a-z]+")

func main() {
	err := db.OpenConnection()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/post", postRevil).Methods("POST")
	r.HandleFunc("/login", loginUser).Methods("POST")
	r.HandleFunc("/register", registerUser).Methods("POST")

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/revil", revilHandler)

	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/register", registerHandler)

	r.HandleFunc("/logout", logoutHandler)

	r.HandleFunc("/user/{username:[a-zA-Z]+}", userHandler)

	r.HandleFunc("/user_taken", userTakenHandler)
	r.HandleFunc("/email_taken", emailTakenHandler)
	r.HandleFunc("/user_valid", isValidUserHandler)

	http.Handle("/", r)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	http.ListenAndServe(":8080", nil)
}

func ensureLoggedIn(writer http.ResponseWriter, request *http.Request) bool {
	loggedIn := isLoggedIn(request)
	if !loggedIn {
		http.Redirect(writer, request, "/login?continue="+request.URL.Path, http.StatusMovedPermanently)
	}

	return loggedIn
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	m := make(map[string]interface{})
	if cont := request.FormValue("continue"); cont != "" {
		m["continue"] = request.FormValue("continue")
	}
	user, _ := getUser(request)

	ShowResponsePage(writer, user, "login", m)
}

func registerHandler(writer http.ResponseWriter, request *http.Request) {
	user, _ := getUser(request)
	ShowResponsePage(writer, user, "register", make(map[string]interface{}))
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	if !ensureLoggedIn(writer, request) {
		return
	}
	user, err := getUser(request)

	revils, err := db.GetAllRevilsInDatabase(*user)
	if err != nil {
		fmt.Println(err)
		revils = make([]data.Revil, 0)
	}

	revilsMap := RevilsAsMap(revils)
	ShowResponsePage(writer, user, "home", revilsMap)
}

func postRevil(writer http.ResponseWriter, request *http.Request) {
	if !ensureLoggedIn(writer, request) {
		return
	}
	user, err := getUser(request)
	if err != nil {
		return
	}

	revilType := request.FormValue("type")
	url := request.FormValue("url")
	title := request.FormValue("title")
	note := request.FormValue("note")
	public := request.FormValue("public") != ""

	id := user.Id.Hex()
	err = db.CreateRevil(id, revilType, url, title, note, public)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}


func revilHandler(writer http.ResponseWriter, request *http.Request) {
	if !ensureLoggedIn(writer, request) {
		return
	}
	user, err := getUser(request)
	if err != nil {
		http.NotFound(writer, request)
	}

	m := make(map[string]interface{})
	if url := request.FormValue("url"); url != "" {
		m["url"] = url
	}
	if revilType := request.FormValue("type"); revilType != "" {
		m[revilType] = true
		m["type"] = revilType
	}
	if title := request.FormValue("title"); title != "" {
		m["title"] = title
	}

	ShowResponsePage(writer, user, "revil", m)
}

func loginUser(writer http.ResponseWriter, request *http.Request) {
	user, canLogin := verifyUser(request)

	if canLogin && user != nil {
		if err := setUser(writer, request, *user); err == nil {
			continueTo := "/" + user.Username
			if val := request.FormValue("continue"); val != "" {
				continueTo = val
			}
			http.Redirect(writer, request, continueTo, http.StatusMovedPermanently)
			return
		} else {
			fmt.Println("Unable to set user", err)
		}
	} else {
		fmt.Println("Unable to login", canLogin, user)
	}
	http.NotFound(writer, request)
}

func userHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	user, err := db.FindUserByName(vars["username"])

	if err != nil || user == nil {
		http.NotFound(writer, request)
		return
	}

	var revils data.Revils
	
	loggedInUser, err := getUser(request)
	if err == nil && loggedInUser.Id == user.Id  {
		revils, err = db.GetAllRevilsInDatabase(*user)
	} else {
		revils, err = db.GetAllPublicRevils(*user)
	}

	if err != nil {
		fmt.Println(err)
		revils = make([]data.Revil, 0)
	}

	values := RevilsAsMap(revils)
	userMap := user.AsMap()

	for key, val := range userMap {
		values[key] = val
	}

	ShowResponsePage(writer, user, "user", values)
	
}

func registerUser(writer http.ResponseWriter, request *http.Request) {
	if isValidRegister(request) {
		username := request.FormValue("username")
		password := request.FormValue("password")
		email := request.FormValue("email")

		user, err := db.CreateUser(username, password, email)
		if err == nil {
			err = setUser(writer, request, user)
			if err == nil {
				http.Redirect(writer, request, "/"+user.Username, http.StatusTemporaryRedirect)
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
}

func isValidRegister(request *http.Request) bool {
	username := request.FormValue("username")
	password := request.FormValue("password")
	verification := request.FormValue("verification")
	email := request.FormValue("email")

	if len(username) < 5 || len(username) > 20 {
		return false
	}
	if !validEmail.MatchString(email) {
		return false
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
		return false
	}

	return true
}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {
	err := logOut(writer, request)
	if err != nil {
		fmt.Println(err)
	}
	ShowResponsePage(writer, nil, "logout", make(map[string]interface{}))
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

func getUser(request *http.Request) (user *data.User, err error) {
	userId, err := getUserId(request)
	if err != nil {
		return
	}

	user, err = db.FindUserById(userId)
	return
}

func ShowResponsePage(writer http.ResponseWriter, user *data.User, name string, data map[string]interface{}) {
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
