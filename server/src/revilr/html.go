package main

import (
	"github.com/hoisie/mustache"
	"revilr/data"
)

var (
	layout         = createPage("resources/html/layout.html")
	navbar         = createPage("resources/html/navbar.html")
	loggedInNavbar = createPage("resources/html/loggedInNavbar.html")
	notFound       = createPage("resources/html/notfound.html")
	internalError  = createPage("resources/html/error.html")
)

var pageMap map[string]*mustache.Template = createPages()

func createPages() map[string]*mustache.Template {
	aMap := make(map[string]*mustache.Template)

	aMap["home"] = createPage("resources/html/display.html")
	aMap["homeNotLoggedIn"] = createPage("resources/html/homeNotLoggedIn.html")
	aMap["login"] = createPage("resources/html/login.html")
	aMap["logout"] = createPage("resources/html/logout.html")
	aMap["user"] = createPage("resources/html/user.html")
	aMap["register"] = createPage("resources/html/register.html")
	aMap["revil"] = createPage("resources/html/revil.html")
	aMap["page"] = createPage("resources/html/display.html")
	aMap["image"] = createPage("resources/html/display.html")
	aMap["selection"] = createPage("resources/html/display.html")
	aMap["emailVerification"] = createPage("resources/html/emailVerification.html")
	aMap["registerSuccessful"] = createPage("resources/html/registerSuccessful.html")

	return aMap
}

func createPage(file string) *mustache.Template {
	tmpl, err := mustache.ParseFile(file)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func RenderWithAdditionalData(page, navbar string, data map[string]interface{}) string {
	data["navbar"] = navbar
	template := pageMap[page]
	if template == nil {
		template = notFound
	}

	return template.RenderInLayout(layout, data)
}

func RenderNotFoundPage() string {
	return notFound.RenderInLayout(layout, make(map[string]interface{}))
}

func RenderInternalErrorPage() string {
	return internalError.RenderInLayout(layout, make(map[string]interface{}))
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

func RenderNavbar(page string, user *data.User) string {
	data := make(map[string]interface{})
	data[page+"Page"] = true
	data["loggedIn"] = user != nil

	var html string

	if user != nil {
		data["user"] = user.AsMap()
		html = loggedInNavbar.Render(data)
	} else {
		html = navbar.Render(data)
	}

	return html
}

func GetEmailVerification(user data.User) string {
	template := pageMap["emailVerification"]
	if template == nil {
		template = notFound
	}
	data := make(map[string]interface{})
	data["verification"] = user.Verification
	data["username"] = user.Username
	return template.Render(data)
}
