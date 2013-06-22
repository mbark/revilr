package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
	"net/url"
)

var index = "templates/index.html"
var layout = "templates/layout.html"

func printAllRevils(revils []revil, writer http.ResponseWriter) {
	data := formatRevilsForOutput(revils)
	html := mustache.RenderFileInLayout(index, layout, data)
	fmt.Fprintf(writer, html)
}

func printAllRevilsOfType(revils []revil, revilType string, writer http.ResponseWriter) {
	data := formatRevilsForOutput(revils)
	data[revilType] = true

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
