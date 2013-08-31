package data

import (
	"labix.org/v2/mgo/bson"
	"net/url"
)

func CreateRevil(revilType, url, title, note string, public bool) (rev Revil) {
	rev = Revil {
		Id: bson.NewObjectId(),
		Created: bson.Now(),
		Type: revilType,
		Url: url,
		Title: title,
		Note: note,
		Public: public,
	}
	
	return
}

func (rev Revil) AsMap() map[string]interface{} {
	data := make(map[string]interface{})

	data["url"] = rev.Url
	data["title"] = rev.Title
	data["note"] = rev.Note
	data["date"] = rev.Created
	data["display-url"] = rev.parseUrl()
	data["type"] = rev.Type
	if(rev.Public) {
		data["public"] = rev.Public		
	}

	return data
}

func (rev Revil) parseUrl() string {
	parsed, err := url.Parse(rev.Url)
	if err != nil {
		return rev.Url
	}

	return parsed.Host
}
