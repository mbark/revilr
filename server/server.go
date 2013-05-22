package main

import (
    "fmt"
    "net/http"
)

func linkHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        fmt.Println("Received link revil!");
        fmt.Println("Url:", r.FormValue("url"));
        fmt.Println("Comment:", r.FormValue("c"));
        fmt.Println();
    }
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        fmt.Println("Received page revil!");
        fmt.Println("Url:", r.FormValue("url"));
        fmt.Println("Comment:", r.FormValue("c"));
        fmt.Println();
    }
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
        fmt.Println("Received image revil!");
        fmt.Println("Url:", r.FormValue("url"));
        fmt.Println("Comment:", r.FormValue("c"));
        fmt.Println();
    }
}

func selectionHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        fmt.Println("Received selection revil!");
        fmt.Println("Url:", r.FormValue("url"));
        fmt.Println("Comment:", r.FormValue("c"));
        fmt.Println();
    }
}

func main() {
    http.HandleFunc("/revilr/link", linkHandler)
    http.HandleFunc("/revilr/page", pageHandler)
    http.HandleFunc("/revilr/image", imageHandler)
    http.HandleFunc("/revilr/selection", selectionHandler)
    http.ListenAndServe(":8080", nil)
}