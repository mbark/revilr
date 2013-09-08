package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"net/smtp"
	"revilr/data"
	"revilr/db"
	"time"
)

var (
	mailAccount  = "noreply.revilr@gmail.com"
	mailPassword = "amenVafan"
)

var (
	ErrInvalidRegistration = errors.New("Invalid parameters for registration")
)

func main() {
	err := db.OpenConnection()
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/post", postRevil).Methods("POST")
	r.HandleFunc("/login", loginUser).Methods("POST")
	r.HandleFunc("/register", registerUser).Methods("POST")

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/revil", revilHandler)

	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/user/{username:[a-zA-Z]+}/verify", verifyRegisterHandler)

	r.HandleFunc("/logout", logoutHandler)

	r.HandleFunc("/user/{username:[a-zA-Z]+}", userHandler)

	r.HandleFunc("/user_taken", userTakenHandler)
	r.HandleFunc("/email_taken", emailTakenHandler)
	r.HandleFunc("/user_valid", isValidUserHandler)

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	http.Handle("/", r)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))

	http.ListenAndServe(":8080", nil)
}

func notFoundHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Unable to find page")

	html := RenderNotFoundPage()
	writer.WriteHeader(404)
	fmt.Fprintf(writer, html)
}

func errorHandler(writer http.ResponseWriter, request *http.Request, err error) {
	fmt.Println("Unexpected error occurred", err)

	html := RenderInternalErrorPage()
	writer.WriteHeader(500)
	fmt.Fprintf(writer, html)
}

func ensureLoggedIn(writer http.ResponseWriter, request *http.Request) (user *data.User, isLoggedIn bool) {
	user, err := getUser(request)
	isLoggedIn = user != nil

	if err != nil {
		err = logOut(writer, request)
		if err != nil {
			errorHandler(writer, request, err)
			return
		}
	}

	if !isLoggedIn {
		http.Redirect(writer, request, "/login?continue="+request.URL.Path, http.StatusMovedPermanently)
	}

	return
}

func getUser(request *http.Request) (user *data.User, err error) {
	userId, err := getUserId(request)
	if err != nil {
		return
	}

	user, err = db.FindUserById(userId)
	return
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
	user, err := getUser(request)

	if user == nil {
		ShowResponsePage(writer, nil, "homeNotLoggedIn", make(map[string]interface{}))
		return
	}

	revils, err := db.GetAllRevils(*user)
	if err != nil {
		fmt.Println("Unable to get all revils", err)
		revils = make([]data.Revil, 0)
	}

	revilsMap := RevilsAsMap(revils)
	ShowResponsePage(writer, user, "home", revilsMap)
}

func postRevil(writer http.ResponseWriter, request *http.Request) {
	user, loggedIn := ensureLoggedIn(writer, request)
	if !loggedIn {
		return
	}

	revilType := request.FormValue("type")
	url := request.FormValue("url")
	title := request.FormValue("title")
	note := request.FormValue("note")
	public := request.FormValue("public") != ""

	id := user.Id.Hex()
	err := db.CreateRevil(id, revilType, url, title, note, public)
	if err != nil {
		fmt.Println("Unable to create revil", err)
	}

	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}

func revilHandler(writer http.ResponseWriter, request *http.Request) {
	user, loggedIn := ensureLoggedIn(writer, request)
	if !loggedIn{
		return
	}

	m := make(map[string]interface{})
	addToMap(request, "url", m)
	addToMap(request, "title", m)
	addToMap(request, "note", m)

	if revilType := request.FormValue("type"); revilType != "" {
		m[revilType] = true
		m["type"] = revilType
	}

	ShowResponsePage(writer, user, "revil", m)
}

func addToMap(request *http.Request, name string, m map[string]interface{}) {
	if val := request.FormValue(name); val != "" {
		m[name] = val
	}
}

func loginUser(writer http.ResponseWriter, request *http.Request) {
	user, canLogin := isValidLogin(request)

	if !canLogin {
		fmt.Println("Unable to login")
		loginHandler(writer, request)
		return
	}

	err := setUser(writer, request, *user)
	if err != nil {
		errorHandler(writer, request, err)
		return
	}

	continueTo := "/user/" + user.Username
	if val := request.FormValue("continue"); val != "" {
		continueTo = val
	}

	http.Redirect(writer, request, continueTo, http.StatusMovedPermanently)
}

func isValidLogin(request *http.Request) (*data.User, bool) {
	username := request.FormValue("username")
	password := request.FormValue("password")

	user, err := db.FindUserByName(username)
	if err == nil && user != nil {
		canLogin := user.PasswordMatches(password)
		canLogin = canLogin && user.Verified
		return user, canLogin
	}
	return nil, false
}


func userHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	user, err := db.FindUserByName(vars["username"])

	if err != nil || user == nil {
		notFoundHandler(writer, request)
		return
	}

	revils := make([]data.Revil, 0)

	loggedInUser, err := getUser(request)

	if loggedInUser != nil && loggedInUser.Id == user.Id {
		revils, err = db.GetAllRevils(*user)
	} else {
		revils, err = db.GetAllPublicRevils(*user)
	}

	values := RevilsAsMap(revils)
	values["user"] = user.AsMap()

	ShowResponsePage(writer, loggedInUser, "user", values)

}

func registerUser(writer http.ResponseWriter, request *http.Request) {
	if !isValidRegister(request) {
		errorHandler(writer, request, ErrInvalidRegistration)
	}

	username := request.FormValue("username")
	password := request.FormValue("password")
	email := request.FormValue("email")

	user, err := db.CreateUser(username, password, email)
	if err != nil {
		errorHandler(writer, request, err)
		return
	}
	
	sendVerificationMail(user)
}

func isValidRegister(request *http.Request) bool {
	username := request.FormValue("username")
	password := request.FormValue("password")
	verification := request.FormValue("verification")
	email := request.FormValue("email")

	if len(username) < 5 || len(username) > 20 {
		return false
	}

	if len(password) < 8 {
		return false
	}

	if password != verification {
		return false
	}

	if db.UserExistsWithName(username) {
		return false
	}

	if db.UserExistsWithEmail(email) {
		return false
	}

	return true
}

func sendVerificationMail(user data.User) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Verification email for revilr\n"
	msg := subject + mime + GetEmailVerification(user)

	auth := smtp.PlainAuth("", mailAccount, mailPassword, "smtp.gmail.com")

	go func(user data.User, msg string) {
		err := smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			mailAccount,
			[]string{user.Email},
			[]byte(msg))

		if err != nil {
			fmt.Println("Unable to send mail", err)
		}
	}(user, msg)
}

func verifyRegisterHandler(writer http.ResponseWriter, request *http.Request) {
	user, err := finishVerification(writer, request)
	if err != nil {
		fmt.Println("Unable to finish verification")
		errorHandler(writer, request, err)
		return
	}

	vals := make(map[string]interface{})
	vals["username"] = user.Username

	ShowResponsePage(writer, user, "registerSuccessful", vals)
}

func finishVerification(writer http.ResponseWriter, request *http.Request) (user *data.User, err error) {
	vars := mux.Vars(request)
	username := vars["username"]
	verification := request.FormValue("verification")

	user, err = db.VerifyUser(username, verification)
	if err != nil {
		return
	}
	err = setUser(writer, request, *user)

	return
}

func logoutHandler(writer http.ResponseWriter, request *http.Request) {
	err := logOut(writer, request)
	if err != nil {
		fmt.Println(err)
	}
	ShowResponsePage(writer, nil, "logout", make(map[string]interface{}))
}

func ShowResponsePage(writer http.ResponseWriter, user *data.User, name string, data map[string]interface{}) {
	navbar := RenderNavbar(name, user)
	html := RenderWithAdditionalData(name, navbar, data)
	fmt.Fprintf(writer, html)
}
