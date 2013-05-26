package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
)

var database *sql.DB

func getRevil(request *http.Request, t string) revil {
	return revil{Type: t, Url: request.FormValue("url"), Comment: request.FormValue("c")}
}

func linkHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		rev := getRevil(request, "link")
		rev.printRevil()
		insertIntoDatabase(rev)
	}
}

func pageHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		rev := getRevil(request, "page")
		rev.printRevil()
		insertIntoDatabase(rev)
	}
}

func imageHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		rev := getRevil(request, "image")
		rev.printRevil()
		insertIntoDatabase(rev)
	}
}

func selectionHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		rev := getRevil(request, "selection")
		rev.printRevil()
		insertIntoDatabase(rev)
	}
}

func getDatabase() (db *sql.DB, err error) {
	dbPath := "./revils.db"

	//check if file exists
	if exists, _ := fileExists(dbPath); exists {
		db, err = sql.Open("sqlite3", "./"+dbPath)
	} else {
		db, err = instantiateDatabase(dbPath)
	}
	return
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func instantiateDatabase(dbPath string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sqls := []string{
		"create table revil (url text not null primary key, type text, comment text)",
	}

	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Printf("%s: %s\n", err, sql)
		}
	}

	return
}

func insertIntoDatabase(rev revil) error {
	stmt, err := database.Prepare("insert into revil(url, type, comment) values(?, ?, ?)")
	if err != nil {
		return err
	}
	stmt.Exec()
	_, err = stmt.Exec(rev.Url, rev.Type, rev.Comment)
	if err != nil {
		return err
	}
	return nil
}

func getAllValuesInDatabase() {
	rows, err := database.Query("select url, type, comment from revil")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		var rtype string
		var comment string
		rows.Scan(&url, &rtype, &comment)
		fmt.Println(url, rtype, comment)
	}
}

func main() {
	db, err := getDatabase()
	if err != nil {
		return
	}

	database = db
	defer database.Close()

	//uncomment to verify it works
	//getAllValuesInDatabase()

	http.HandleFunc("/revilr/link", linkHandler)
	http.HandleFunc("/revilr/page", pageHandler)
	http.HandleFunc("/revilr/image", imageHandler)
	http.HandleFunc("/revilr/selection", selectionHandler)
	http.ListenAndServe(":8080", nil)
}
