package main

import (
	"html/template"
	"net/http"
	"regexp"
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
	rev := getAllRevilsInDatabase()

	t, _ := template.ParseFiles("showall.html")
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
	http.ListenAndServe(":8080", nil)
}
