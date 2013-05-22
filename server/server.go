package main

import (
    "fmt"
    "net/http"
)

func linkHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    fmt.Println("Reviled link",values["val"], "with comment", values["c"])
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    fmt.Println("Reviled page",values["val"], "with comment", values["c"])
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
    fmt.Println("Reviled image",values["val"], "with comment", values["c"])
}

func selectionHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    fmt.Println("Reviled selection",values["val"], "with comment", values["c"])
}

func main() {
    http.HandleFunc("/revil/link", linkHandler)
    http.HandleFunc("/revil/page", pageHandler)
    http.HandleFunc("/revil/image", imageHandler)
    http.HandleFunc("/revil/selection", selectionHandler)
    http.ListenAndServe(":8080", nil)
}