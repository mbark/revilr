package data

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"net/url"
)

func CreateRevil(revilType, url, title, note string) (rev Revil) {
	rev.Id = bson.NewObjectId()
	rev.Created = bson.Now()
	rev.Type = revilType
	rev.Url = url
	rev.Title = title
	rev.Note = note
	return
}

func (r Revil) toString() string {
	var result string

	result += "{"
	result += " type: " + r.Type
	result += ", url: " + r.Url
	result += ", title: " + r.Title
	result += ", note: " + r.Note
	result += " }"

	return result
}

func (r Revil) Print() {
	fmt.Println(r.toString())
}

func (rev Revil) AsMap() map[string]interface{} {
	data := make(map[string]interface{})

	data["url"] = rev.Url
	data["title"] = rev.Title
	data["note"] = rev.Note
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
