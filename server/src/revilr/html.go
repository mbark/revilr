package main

import (
	"github.com/hoisie/mustache"
	"revilr/data"
)

var (
	layout   = createPage("resources/html/layout.html")
	navbar   = createPage("resources/html/navbar.html")
	notFound = createPage("resources/html/notfound.html")
)

var pageMap map[string]*mustache.Template = createPages()

func createPages() map[string]*mustache.Template {
	aMap := make(map[string]*mustache.Template)

	aMap["home"] = createPage("resources/html/display.html")
	aMap["login"] = createPage("resources/html/login.html")
	aMap["logout"] = createPage("resources/html/logout.html")
	aMap["user"] = createPage("resources/html/user.html")
	aMap["register"] = createPage("resources/html/register.html")
	aMap["revil"] = createPage("resources/html/revil.html")
	aMap["page"] = createPage("resources/html/display.html")
	aMap["image"] = createPage("resources/html/display.html")
	aMap["selection"] = createPage("resources/html/display.html")

	return aMap
}

func createPage(file string) *mustache.Template {
	tmpl, err := mustache.ParseFile(file)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func Render(page string, user *data.User) string {
	return RenderWithAdditionalData(page, user, make(map[string]interface{}))
}

func RenderWithAdditionalData(page string, user *data.User, data map[string]interface{}) string {
	data["navbar"] = getNavbar(page, user)
	template := pageMap[page]
	if template == nil {
		template = notFound
	}

	return template.RenderInLayout(layout, data)
}

func RevilsAsMap(revils []data.Revil) map[string]interface{} {
	revilsMap := make([]map[string]interface{}, len(revils))

	for key, rev := range revils {
		dataType := make(map[string]interface{})
		dataType[rev.Type] = rev.AsMap()

		revilsMap[key] = dataType
	}

	values := make(map[string]interface{})
	values["revils"] = revilsMap
	return values
}

func getNavbar(page string, user *data.User) string {
	data := make(map[string]interface{})
	data[page] = true
	data["loggedIn"] = user != nil

	if user != nil {
		data["username"] = user.Username
		data["emailHash"] = user.EmailHash()
	}

	return navbar.Render(data)
}
