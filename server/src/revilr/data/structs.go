package data

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	Revils []Revil

	Revil struct {
		Id      bson.ObjectId `json:"id"      bson:"_id"`
		UserId  bson.ObjectId `json:"uid"     bson:"uid"`
		Type    string        `json:"type"    bson:"type"`
		Url     string        `json:"url"     bson:"url"`
		Comment string        `json:"comment" bson:"comment"`
		Created time.Time     `json:"created" bson:"created"`
	}

	User struct {
		Id       bson.ObjectId `json:"id"       bson:"_id"`
		Username string        `json:"username" bson:"username"`
		Password []byte        `json:"password" bson:"password"`
		Created  time.Time     `json:"created"  bson:"created"`
	}
)
