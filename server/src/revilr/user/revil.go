package user

import (
	"fmt"
	"net/url"
)

type Revil struct {
	Type    string
	Url     string
	Comment string
	Date    string
}

func (r Revil) toString() string {
	var result string

	result += "{"
	result += " type: " + r.Type
	result += ", url: " + r.Url
	result += ", comment: \"" + r.Comment + "\""
	result += " }"

	return result
}

func (r Revil) PrintRevil() {
	fmt.Println(r.toString())
}

func (rev Revil) AsMap() map[string]interface{} {
	data := make(map[string]interface{})

	data["url"] = rev.Url
	data["comment"] = rev.Comment
	data["date"] = rev.Date
	data["display-url"] = parseUrl(rev)
	data["type"] = rev.Type

	return data
}


func parseUrl(rev Revil) string {
	parsed, err := url.Parse(rev.Url)
	if err != nil {
		fmt.Println(err)
		return rev.Url
	}

	return parsed.Host
}
