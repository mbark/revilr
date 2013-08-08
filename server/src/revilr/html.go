package main

import (
	"fmt"
	"github.com/hoisie/mustache"
	"net/http"
	"revilr/data"
)

var (
	layout    = parseFile("resources/html/layout.html")
	navbar    = parseFile("resources/html/navbar.html")
	display   = parseFile("resources/html/display.html")
	loginTmpl = parseFile("resources/html/login.html")
	logout    = parseFile("resources/html/logout.html")
	userUrl   = parseFile("resources/html/user.html")
	register  = parseFile("resources/html/register.html")
	revil     = parseFile("resources/html/revil.html")
)

func parseFile(file string) *mustache.Template {
	tmpl, err := mustache.ParseFile(file)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func DisplayRevils(revils []data.Revil, revilType string, writer http.ResponseWriter, request *http.Request) {
	data := formatRevilsForOutput(revils)
	data["navbar"] = getNavbar(revilType, request)
	html := display.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayLogin(writer http.ResponseWriter, request *http.Request) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("login", request)
	html := loginTmpl.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayLogout(writer http.ResponseWriter, request *http.Request) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("logout", request)
	html := logout.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayUser(writer http.ResponseWriter, request *http.Request, user data.User) {
	data := user.AsMap()
	data["navbar"] = getNavbar("user", request)
	html := userUrl.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayRegister(writer http.ResponseWriter, request *http.Request) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("register", request)
	html := register.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func DisplayRevil(writer http.ResponseWriter, request *http.Request) {
	data := make(map[string]interface{})
	data["navbar"] = getNavbar("revil", request)
	html := revil.RenderInLayout(layout, data)

	fmt.Fprintf(writer, html)
}

func formatRevilsForOutput(revils []data.Revil) map[string]interface{} {
	values := make(map[string]interface{})
	values["revils"] = getListOfRevilMaps(revils)
	return values
}

func getListOfRevilMaps(revils []data.Revil) []map[string]interface{} {
	revilMaps := make([]map[string]interface{}, len(revils))

	for key, rev := range revils {
		revilMaps[key] = getMapForRevil(rev)
	}

	return revilMaps
}

func getMapForRevil(rev data.Revil) map[string]interface{} {
	dataType := make(map[string]interface{})
	dataType[rev.Type] = rev.AsMap()

	return dataType
}

func getNavbar(revilType string, request *http.Request) string {
	data := make(map[string]interface{})
	data[revilType] = true
	data["loggedIn"] = false

	if user, ok := getUser(request); ok && isLoggedIn(request) {
		data["username"] = user.Username
		data["loggedIn"] = true
		data["emailHash"] = user.EmailHash()
	}

	return navbar.Render(data)
}
