package main

import (
    "fmt"
    "net/http"
)

func linkHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    fmt.Println("Revild link",values["val"], "with comment", values["c"])
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    fmt.Println("Revild page",values["val"], "with comment", values["c"])
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
    fmt.Println("Revild image",values["val"], "with comment", values["c"])
}

func selectionHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    fmt.Println("Revild selection",values["val"], "with comment", values["c"])
}

func main() {
    http.HandleFunc("/revilr/link", linkHandler)
    http.HandleFunc("/revilr/page", pageHandler)
    http.HandleFunc("/revilr/image", imageHandler)
    http.HandleFunc("/revilr/selection", selectionHandler)
    http.ListenAndServe(":8080", nil)
}