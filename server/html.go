package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
	"net/url"
)

var index = "templates/index.html"
var layout = "templates/layout.html"
var navbar = "templates/navbar.html"

func printAllRevils(revils []revil, writer http.ResponseWriter) {
	data := formatRevilsForOutput(revils)
	data["navbar"] = getNavbar()
	html := mustache.RenderFileInLayout(index, layout, data)
	fmt.Fprintf(writer, html)
}

func printAllRevilsOfType(revils []revil, revilType string, writer http.ResponseWriter) {
	data := formatRevilsForOutput(revils)
	data["navbar"] = getNavbarForType(revilType)

	htmlFile := "templates/" + revilType + ".html"
	html := mustache.RenderFileInLayout(htmlFile, layout, data)
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
	data := make(map[string]interface{})

	data["url"] = rev.Url
	data["comment"] = rev.Comment
	data["date"] = rev.Date
	data["display-url"] = parseUrl(rev)

	dataType := make(map[string]interface{})
	dataType[rev.Type] = data

	return dataType
}

func parseUrl(rev revil) string {
	parsed, err := url.Parse(rev.Url)
	if err != nil {
		fmt.Println(err)
		return rev.Url
	}

	return parsed.Host
}

func getNavbar() string {
	data := make(map[string]interface{})
	html := mustache.RenderFile(navbar, data)

	return html
}

func getNavbarForType(revilType string) string {
	data := make(map[string]interface{})
	data[revilType] = true
	html := mustache.RenderFile(navbar, data)

	return html
}