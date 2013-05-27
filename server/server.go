package main

import (
	"net/http"
    "regexp"
)

const lenPath = len("/revilr/")
var typeValidator = regexp.MustCompile("^(link|page|image|selection)$")

func getRevil(request *http.Request, t string) revil {
	return revil{Type: t, Url: request.FormValue("url"), Comment: request.FormValue("c")}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		postHandler(w, r)
	} else if r.Method == "GET" {

	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	rType := r.URL.Path[lenPath:]
	if !typeValidator.MatchString(rType) {
		http.NotFound(w, r)
		return
	}
	rev := getRevil(r, rType)
	rev.printRevil()
	insertIntoDatabase(rev)
}

func main() {
	db, err := getDatabase()
	if err != nil {
		return
	}

	database = db
	defer database.Close()

	//uncomment to verify it works
	getAllValuesInDatabase()

	http.HandleFunc("/revilr/", httpHandler)
	http.ListenAndServe(":8080", nil)
}
