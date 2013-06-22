package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
	"net/url"
)

var layout = "templates/layout.html"
var navbar = "templates/navbar.html"

func printAllRevils(revils []revil, revilType string, writer http.ResponseWriter) {
	data := formatRevilsForOutput(revils, revilType)
	data["navbar"] = getNavbar(revilType)

	htmlFile := "templates/" + revilType + ".html"
	html := mustache.RenderFileInLayout(htmlFile, layout, data)

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
	html := mustache.RenderFile(navbar, data)

	return html
}
