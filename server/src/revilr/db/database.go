package db

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"revilr/data"
)

const (
	name    = "revilr"
	url     = "127.0.0.1"
	usersC  = "users"
	revilsC = "revils"
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

func CreateRevil(userId, revilType, url, title, note string) error {
	rev := data.CreateRevil(revilType, url, title, note)
	rev.UserId = bson.ObjectIdHex(userId)

	collection := database.C(revilsC)
	_, err := collection.UpsertId(rev.Id, rev)
	return err
}

func GetAllRevilsInDatabase(user data.User) (revils data.Revils, err error) {
	collection := database.C(revilsC)
	err = collection.Find(bson.M{"uid": user.Id}).Sort("-created").All(&revils)
	return
}

func GetRevilsOfType(rtype string, user data.User) (revils data.Revils, err error) {
	collection := database.C(revilsC)
	err = collection.Find(bson.M{"type": rtype, "uid": user.Id}).Sort("-created").All(&revils)
	return
}

func FindUserById(userId string) (user *data.User, err error) {
	collection := database.C(usersC)
	id := bson.ObjectIdHex(userId)
	err = collection.Find(bson.M{"_id": id}).One(&user)
	return
}

func FindUserByName(username string) (user *data.User, err error) {
	collection := database.C(usersC)
	err = collection.Find(bson.M{"username": username}).One(&user)
	return
}

func CreateUser(username, password string) (user data.User, err error) {
	user, err = data.CreateUser(username, password)
	if err != nil {
		return
	}

	collection := database.C(usersC)
	_, err = collection.UpsertId(user.Id, user)
	return
}
