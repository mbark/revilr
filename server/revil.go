package main

import (
	"fmt"
)

type revil struct {
	Type    string
	Url     string
	Comment string
	Date    string
}

func (r revil) toString() string {
	var result string

	result += "{"
	result += " type: " + r.Type
	result += ", url: " + r.Url
	result += ", comment: \"" + r.Comment + "\""
	result += " }"

	return result
}

func (r revil) printRevil() {
	fmt.Println(r.toString())
}
