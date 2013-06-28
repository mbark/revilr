package main

import (
	"net/http"
	"regexp"
	"fmt"
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
	displayRevils(revils, revilType, writer)
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	revils := getAllRevilsInDatabase()
	displayRevils(revils, "all", writer)
}

func userHandler(writer http.ResponseWriter, request *http.Request) {
	username, password := parseUser(request)
	user, err := getUser(username)
	if err != nil {
		fmt.Println(err)
		return
	}

	if user.Username == "" {
		createUser(username, password)
		fmt.Println("Made user", username, "with password", password)
	} else {
		fmt.Println("Found matching user", user)
	}
}

func parseUser(request *http.Request) (username, password string) {
	username = request.FormValue("user")
	password = request.FormValue("password")
	return
}