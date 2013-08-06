package db

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"revilr/data"
	"time"
)

const (
	name = "revilr"
	url  = "127.0.0.1"
)

var database *mgo.Database

func OpenConnection() (err error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return
	}

	database = session.DB(name)
	return
}

func InsertIntoDatabase(rev data.Revil, user data.User) {
	collection := database.C("revils")
	rev.UserId = user.Id
	rev.Id = bson.NewObjectId()
	rev.Created = time.Now()
	collection.UpsertId(rev.Id, rev)
	return
}

func GetAllRevilsInDatabase(usr data.User) (revils data.Revils, err error) {
	collection := database.C("revils")
	err = collection.Find(bson.M{"_uid": usr.Id}).All(&revils)
	return
}

func GetRevilsOfType(rtype string, usr data.User) (revils data.Revils, err error) {
	collection := database.C("revils")
	err = collection.Find(bson.M{"type": rtype, "_uid": usr.Id}).All(&revils)
	return
}

func FindUser(username string) (user *data.User, err error) {
	collection := database.C("users")
	err = collection.Find(bson.M{"username": username}).One(&user)
	return
}

func CreateUser(user data.User) {
	collection := database.C("users")
	user.Id = bson.NewObjectId()
	user.Created = time.Now()
	collection.UpsertId(user.Id, user)
	return
}
