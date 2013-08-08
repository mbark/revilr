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
		Title	string        `json:"title"   bson:"title"`
		Note 	string        `json:"note"    bson:"note"`
		Created time.Time     `json:"created" bson:"created"`
	}

	User struct {
		Id       bson.ObjectId `json:"id"       bson:"_id"`
		Username string        `json:"username" bson:"username"`
		Password []byte        `json:"password" bson:"password"`
		Email    string        `json="email"    bson:"email"`
		Created  time.Time     `json:"created"  bson:"created"`
	}
)
