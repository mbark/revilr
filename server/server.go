package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

const lenPath = len("/revilr/")

var validTypes = regexp.MustCompile("^(page|image|selection)$")

func getRevil(request *http.Request, t string) revil {
	return revil{Type: t, Url: request.FormValue("url"), Comment: request.FormValue("c")}
}

func main() {
	db, err := getDatabase()
	if err != nil {
		return
	}

	database = db
	defer database.Close()

	http.HandleFunc("/revilr/", httpHandler)
	http.HandleFunc("/revilr", indexHandler)
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

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	rev := getAllRevilsInDatabase()
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(writer, rev)
}

func postHandler(request *http.Request, revilType string) {
	rev := getRevil(request, revilType)
	rev.printRevil()
	insertIntoDatabase(rev)
}

func getHandler(writer http.ResponseWriter, request *http.Request, revilType string) {
	rev := getRevilOfType(revilType)

	htmlFile := "templates/" + revilType + ".html"
	t, err := template.ParseFiles(htmlFile)

	if err != nil {
		fmt.Println(err)
		http.NotFound(writer, request)
	}
	t.Execute(writer, rev)
}
