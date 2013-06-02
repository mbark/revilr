package main

import (
	"fmt"
)

type revil struct {
	Type    string
	Url     string
	Comment string
}

func (r revil) toString() string {
	result := r.Type + " revil"
	result += "\nUrl: " + r.Url
	if len(r.Comment) != 0 {
		result += "\nComment: " + r.Comment
	}
	return result
}

func (r revil) printRevil() {
	fmt.Println("--------------")
	fmt.Println(r.toString())
	fmt.Println("--------------")
	fmt.Println()
}
