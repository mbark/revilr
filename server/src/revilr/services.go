package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"revilr/db"
)

func emailTakenHandler(writer http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")
	user, _ := db.FindUserByEmail(email)
	isTaken := user != nil

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, Response{"isTaken": isTaken})
}

func userTakenHandler(writer http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	isTaken := db.UserExistsWithName(username)

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, Response{"isTaken": isTaken})
}

func isValidUserHandler(writer http.ResponseWriter, request *http.Request) {
	_, canLogin := isValidLogin(request)

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprint(writer, Response{"isValid": canLogin})
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}
