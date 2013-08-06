package data

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type (
	Revils []Revil

	Revil struct {
		Id      bson.ObjectId `json:"rid"	bson:"_rid"`
		UserId  bson.ObjectId `json:"uid" bson:"_uid"`
		Type    string        `json:"type" bson:"type"`
		Url     string        `json:"url" bson:"url"`
		Comment string        `json:"comment" bson:"comment"`
		Created time.Time     `json:"date" bson:"date"`
	}

	User struct {
		Id       bson.ObjectId `json:"uid" bson:"_uid"`
		Username string        `json:"username" bson:"username"`
		Password []byte        `json:"password" bson:"password"`
		Created  time.Time     `json:"date" bson:"date"`
	}
)
