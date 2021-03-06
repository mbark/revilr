package db

import (
	"errors"
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

var (
	ErrVerificationDoesNotMatch = errors.New("Verification does not match expected value")
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

func CreateRevil(userId, revilType, url, title, note string, public bool) error {
	rev := data.CreateRevil(revilType, url, title, note, public)
	rev.UserId = bson.ObjectIdHex(userId)

	collection := database.C(revilsC)
	_, err := collection.UpsertId(rev.Id, rev)
	return err
}

func DeleteRevil(id string) error {
	bsonId := bson.ObjectIdHex(id)
	collection := database.C(revilsC)
	err := collection.RemoveId(bsonId)
	return err
}

func GetAllRevils(user data.User) (revils data.Revils, err error) {
	collection := database.C(revilsC)
	err = collection.Find(bson.M{"uid": user.Id}).Sort("-created").All(&revils)
	return
}

func GetAllPublicRevils(user data.User) (revils data.Revils, err error) {
	collection := database.C(revilsC)
	err = collection.Find(bson.M{"public": true, "uid": user.Id}).Sort("-created").All(&revils)
	return
}

func UserExistsWithName(name string) bool {
	user, _ := FindUserByName(name)
	return user != nil
}

func UserExistsWithEmail(email string) bool {
	user, _ := FindUserByEmail(email)
	return user != nil
}

func FindUserById(userId string) (user *data.User, err error) {
	collection := database.C(usersC)
	id := bson.ObjectIdHex(userId)
	err = collection.FindId(id).One(&user)
	return
}

func FindUserByName(username string) (user *data.User, err error) {
	collection := database.C(usersC)
	err = collection.Find(bson.M{"username": username}).One(&user)
	return
}

func FindUserByEmail(email string) (user *data.User, err error) {
	collection := database.C(usersC)
	err = collection.Find(bson.M{"email": email}).One(&user)
	return
}

func CreateUser(username, password, email string) (user data.User, err error) {
	user, err = data.CreateUser(username, password, email)
	if err != nil {
		return
	}

	collection := database.C(usersC)
	_, err = collection.UpsertId(user.Id, user)
	return
}

func VerifyUser(username, verification string) (user *data.User, err error) {
	collection := database.C(usersC)
	err = collection.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return
	}

	if user.Verification != verification {
		err = ErrVerificationDoesNotMatch
	} else {
		change := bson.M{"$set": bson.M{"verified": true}, "$unset": bson.M{"verification": ""}}
		collection.Update(bson.M{"_id": user.Id}, change)

		user, err = FindUserById(user.Id.Hex())
	}

	return
}
