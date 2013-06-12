package main

import (
	"html/template"
	"net/http"
	"regexp"
	"fmt"
)

const lenPath = len("/revilr/")

var postTypeValidator = regexp.MustCompile("^(link|page|image|selection)$")
var getTypeValidator = regexp.MustCompile("^(link|page|image|selection|)$")

func getRevil(request *http.Request, t string) revil {
	return revil{Type: t, Url: request.FormValue("url"), Comment: request.FormValue("c")}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postHandler(w, r)
	} else if r.Method == "GET" {
		getHandler(w, r)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	rType := r.URL.Path[lenPath:]
	if !postTypeValidator.MatchString(rType) {
		http.NotFound(w, r)
		return
	}
	rev := getRevil(r, rType)
	rev.printRevil()
	insertIntoDatabase(rev)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	rType := r.URL.Path[lenPath:]
	if !getTypeValidator.MatchString(rType) {
		http.NotFound(w, r)
		return
	}
	fmt.Println(rType)
	rev := getRevilOfType(rType)

	htmlFile := "templates/" + rType + ".html"
	t, err := template.ParseFiles(htmlFile)

	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
	}
	t.Execute(w, rev)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rev := getAllRevilsInDatabase()

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, rev)
}

func main() {
	db, err := getDatabase()
	if err != nil {
		return
	}

	database = db
	defer database.Close()

	//uncomment to verify it works
	printAllRevilsInDatabase()

	http.HandleFunc("/revilr/", httpHandler)
	http.HandleFunc("/revilr", indexHandler)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("templates/resources"))))
	http.ListenAndServe(":8080", nil)
}
