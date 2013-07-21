package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
	"revilr/user"
)

var layout = parseFile("resources/html/layout.html")
var navbar = parseFile("resources/html/navbar.html")
var display = parseFile("resources/html/display.html")
var login = parseFile("resources/html/login.html")
var logout = parseFile("resources/html/logout.html")
var userUrl = parseFile("resources/html/user.html")
var register = parseFile("resources/html/register.html")

func parseFile(file string) *mustache.Template {
	tmpl, err := mustache.ParseFile(file)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func DisplayRevils(revils []user.Revil, revilType string, writer http.ResponseWriter, request *http.Request) {
	data := formatRevilsForOutput(revils, revilType)
	data["navbar"] = getNavbar(revilType, request)
	html := display.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayLogin(writer http.ResponseWriter, success string, request *http.Request) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("login", request)
	data[success] = true
	html := login.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayLogout(writer http.ResponseWriter, loggedIn bool, request *http.Request) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("logout", request)
	data["loggedIn"] = loggedIn
	html := logout.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayUser(writer http.ResponseWriter, request *http.Request) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("user", request)
	html := userUrl.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayRegister(writer http.ResponseWriter, success string, request *http.Request) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("register", request)
	data[success] = true
	html := register.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func formatRevilsForOutput(revils []user.Revil, revilType string) map[string]interface{} {
	values := make(map[string]interface{})
	values["revils"] = getListOfRevilMaps(revils)
	return values
}

func getListOfRevilMaps(revils []user.Revil) []map[string]interface{} {
	revilMaps := make([]map[string]interface{}, len(revils))

	for key, rev := range revils {
		revilMaps[key] = getMapForRevil(rev)
	}

	return revilMaps
}

func getMapForRevil(rev user.Revil) map[string]interface{} {
	dataType := make(map[string]interface{})
	dataType[rev.Type] = rev.AsMap()

	return dataType
}

func getNavbar(revilType string, request *http.Request) string {
	loggedIn, username := getUsername(request)
	data := make(map[string]interface{})
	data[revilType] = true
	data["loggedIn"] = loggedIn
	if loggedIn {
		data["username"] = username
	}
	html := navbar.Render(data)

	return html
}
