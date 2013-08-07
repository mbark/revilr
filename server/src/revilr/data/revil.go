package data

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"net/url"
)

func CreateRevil(revilType, url, comment string) (rev Revil) {
	rev = Revil{Id: bson.NewObjectId(), Type: revilType, Url: url, Comment: comment, Created: bson.Now()}
	return
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

func (r Revil) Print() {
	fmt.Println(r.toString())
}

func (rev Revil) AsMap() map[string]interface{} {
	data := make(map[string]interface{})

	data["url"] = rev.Url
	data["comment"] = rev.Comment
	data["date"] = rev.Created
	data["display-url"] = rev.parseUrl()
	data["type"] = rev.Type

	return data
}

func (rev Revil) parseUrl() string {
	parsed, err := url.Parse(rev.Url)
	if err != nil {
		return rev.Url
	}

	return parsed.Host
}
