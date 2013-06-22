package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"github.com/hoisie/mustache"
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

func postHandler(request *http.Request, revilType string) {
	rev := getRevil(request, revilType)
	rev.printRevil()
	insertIntoDatabase(rev)
}

func getRevil(request *http.Request, t string) revil {
	return revil{Type: t, Url: request.FormValue("url"), Comment: request.FormValue("c")}
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	revils := getAllRevilsInDatabase()
	htmlFile := "templates/index.html"
	printAllRevils(revils, htmlFile, writer)
}

func getHandler(writer http.ResponseWriter, request *http.Request, revilType string) {
	revils := getRevilsOfType(revilType)
	htmlFile := "templates/" + revilType + ".html"
	printAllRevils(revils, htmlFile, writer)
}

func printAllRevils(revils []revil, htmlFile string, writer http.ResponseWriter) {
	data := formatRevilsForOutput(revils)
	html := mustache.RenderFile(htmlFile, data)
	fmt.Fprintf(writer, html)
}

func formatRevilsForOutput(revils []revil) map[string]interface{} {
	values := make(map[string]interface{})
		
	values["revils"] = getListOfRevilMaps(revils)

	return values
}

func getListOfRevilMaps(revils []revil) []map[string]interface{} {
	revilMaps := make([]map[string]interface{}, len(revils))

	for key, rev := range revils {
		revilMaps[key] = getMapForRevil(rev)
	}

	return revilMaps
}

func getMapForRevil(rev revil) map[string]interface{} {
	values := make(map[string]interface{})

	values["type"] = rev.Type
	values["url"] = rev.Url
	values["comment"] = rev.Comment
	values["date"] = rev.Date
	values["display-url"] = parseUrl(rev)
	
	return values
}

func parseUrl(rev revil) string {
	parsed, err := url.Parse(rev.Url)
	if err != nil {
		fmt.Println(err)
		return rev.Url
	}

	return parsed.Host
}