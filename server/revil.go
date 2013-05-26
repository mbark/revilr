package main

import (
	"fmt"
)

type revil struct {
	Type    string
	Url     string
	Comment string
}

func (r revil) getString() string {
	result := "Url: " + r.Url
	if len(r.Comment) != 0 {
		result += "\nComment: " + r.Comment
	}
	return result
}

func (r revil) printRevil() {
	fmt.Println("--------------")
	fmt.Println(r.Type, "revil:")
	fmt.Println(r.getString())
	fmt.Println("--------------")
	fmt.Println()
}
