package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
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
	revils := getAllRevilsInDatabase()
	for key,rev := range revils {
		revils[key].DisplayUrl = parseUrl(rev)
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(writer, revils)
}

func postHandler(request *http.Request, revilType string) {
	rev := getRevil(request, revilType)
	rev.printRevil()
	insertIntoDatabase(rev)
}

func getRevil(request *http.Request, t string) revil {
	return revil{Type: t, Url: request.FormValue("url"), Comment: request.FormValue("c")}
}


func getHandler(writer http.ResponseWriter, request *http.Request, revilType string) {
	revils := getRevilsOfType(revilType)
	for key,rev := range revils {
		revils[key].DisplayUrl = parseUrl(rev)
	}

	htmlFile := "templates/" + revilType + ".html"
	t, err := template.ParseFiles(htmlFile)

	if err != nil {
		fmt.Println(err)
		http.NotFound(writer, request)
	}
	t.Execute(writer, revils)
}

func parseUrl(rev revil) string {
	parsed, err := url.Parse(rev.Url)
	if err != nil {
		fmt.Println(err)
		return rev.Url
	}
	fmt.Println(parsed.Host)

	return parsed.Host
}