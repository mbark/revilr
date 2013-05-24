package main

import (
    "encoding/json"
    "io"
    "os"
)

var file *os.File

func initializeDb() {
	filename := "revils.db"

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil { panic(err) }
	file = f;
}

func toJson(r revil) []byte {
	b, err := json.Marshal(r)
	if err != nil { panic(err) }

	return b
}

//Untested
func toRevil(b []byte) revil {
	var r revil

	err := json.Unmarshal(b, &r)
	if err != nil { panic(err) }

	return r
}

func save(rev revil) {
    b := toJson(rev)

    _, err := io.WriteString(file, string(b) + "\n")
	if err != nil { panic(err) }
}