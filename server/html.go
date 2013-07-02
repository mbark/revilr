package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
	"net/url"
)

var layout = parseFile("templates/layout.html")
var navbar = parseFile("templates/navbar.html")
var display = parseFile("templates/display.html")
var login = parseFile("templates/login.html")

func parseFile(file string) *mustache.Template {
	tmpl, err := mustache.ParseFile(file)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func DisplayRevils(revils []revil, revilType string, writer http.ResponseWriter) {
	data := formatRevilsForOutput(revils, revilType)
	data["navbar"] = getNavbar(revilType)
	html := display.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayLogin(writer http.ResponseWriter) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("login")
	html := login.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func formatRevilsForOutput(revils []revil, revilType string) map[string]interface{} {
	values := make(map[string]interface{})
	values["revils"] = getListOfRevilMaps(revils)
	values["navbar"] = getNavbar(revilType)
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
	dataType := make(map[string]interface{})
	dataType[rev.Type] = rev.asMap()

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

func getNavbar(revilType string) string {
	data := make(map[string]interface{})
	data[revilType] = true
	html := navbar.Render(data)

	return html
}
