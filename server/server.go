package main

import (
    "fmt"
    "net/http"
)

func revilHandler(w http.ResponseWriter, r *http.Request) {
    values := r.URL.Query()
    fmt.Println("Reviled",values["page"])
}

func main() {
    http.HandleFunc("/revil", revilHandler)
    http.ListenAndServe(":8080", nil)
}