package main

import (
    "fmt"
    "net/http"
)

type revil struct {
    rType string
    rUrl string
    rComment string
}

func getString(r revil) string {
    result := "Url: " + r.rUrl
    if len(r.rComment) != 0 {
        result += "\nComment: " + r.rComment;
    }
    return result
}

func printRequest(request *http.Request, t string) {
    rev := revil{rType:t, rUrl:request.FormValue("url"), rComment:request.FormValue("c")}

    fmt.Println("Received", rev.rType, "revil!")
    fmt.Println(getString(rev))
    fmt.Println()
}

func linkHandler(w http.ResponseWriter, request *http.Request) {
    if request.Method == "POST" {
        printRequest(request, "link")
    }
}

func pageHandler(w http.ResponseWriter, request *http.Request) {
    if request.Method == "POST" {
        printRequest(request, "page")
    }
}

func imageHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
        printRequest(request, "image")
    }
}

func selectionHandler(w http.ResponseWriter, request *http.Request) {
    if request.Method == "POST" {
        printRequest(request, "selection")
    }
}

func main() {
    http.HandleFunc("/revilr/link", linkHandler)
    http.HandleFunc("/revilr/page", pageHandler)
    http.HandleFunc("/revilr/image", imageHandler)
    http.HandleFunc("/revilr/selection", selectionHandler)
    http.ListenAndServe(":8080", nil)
}