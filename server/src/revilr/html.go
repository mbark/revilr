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

func DisplayRevils(revils []user.Revil, revilType string, writer http.ResponseWriter) {
	data := formatRevilsForOutput(revils, revilType)
	data["navbar"] = getNavbar(revilType)
	html := display.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayLogin(writer http.ResponseWriter, success string) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("login")
	data[success] = true
	html := login.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayLogout(writer http.ResponseWriter, isLoggedOut string) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("logout")
	data[isLoggedOut] = true
	html := logout.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayUser(writer http.ResponseWriter, username string) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbarUser("user", username)
	data["username"] = username
	html := userUrl.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayRegister(writer http.ResponseWriter, success string) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("register")
	data[success] = true
	html := register.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func formatRevilsForOutput(revils []user.Revil, revilType string) map[string]interface{} {
	values := make(map[string]interface{})
	values["revils"] = getListOfRevilMaps(revils)
	values["navbar"] = getNavbar(revilType)
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

func getNavbar(revilType string) string {
	data := make(map[string]interface{})
	data[revilType] = true
	data["loggedIn"] = false
	html := navbar.Render(data)

	return html
}

func getNavbarUser(revilType string, username string) string {
	data := make(map[string]interface{})
	data[revilType] = true
	data["username"] = username
	data["loggedIn"] = true
	html := navbar.Render(data)

	return html
}
