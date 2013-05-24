package main

import (
    "net/http"
)

func getRevil(request *http.Request, t string) revil {
    return revil{Type:t, Url:request.FormValue("url"), Comment:request.FormValue("c")};
}

func linkHandler(w http.ResponseWriter, request *http.Request) {
    if request.Method == "POST" {
        rev := getRevil(request, "link");
        rev.printRevil();
        save(rev)
    }
}

func pageHandler(w http.ResponseWriter, request *http.Request) {
    if request.Method == "POST" {
        rev := getRevil(request, "page");
        rev.printRevil();
        save(rev)
    }
}

func imageHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
        rev := getRevil(request, "image");
        rev.printRevil();
        save(rev)
    }
}

func selectionHandler(w http.ResponseWriter, request *http.Request) {
    if request.Method == "POST" {
        rev := getRevil(request, "selection");
        rev.printRevil();
        save(rev)
    }
}

func main() {
    initializeDb();
    http.HandleFunc("/revilr/link", linkHandler)
    http.HandleFunc("/revilr/page", pageHandler)
    http.HandleFunc("/revilr/image", imageHandler)
    http.HandleFunc("/revilr/selection", selectionHandler)
    http.ListenAndServe(":8080", nil)
}