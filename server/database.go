package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

var database *sql.DB

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

func printAllRevilsInDatabase() {
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
		revil{Type: rtype, Url: url, Comment: comment}.printRevil()
	}
}

func getAllRevilsInDatabase() []revil {
	rows, err := database.Query("select url, type, comment from revil order by ROWID desc")
	if err != nil {
		fmt.Println("Error:", err)
		return make([]revil, 0)
	}
	defer rows.Close()

<<<<<<< HEAD
=======
	return rowsToRevils(rows)
}

func getRevilOfType(rtype string) []revil {
	rows, err := database.Query("select url, type, comment from revil WHERE type=?", rtype)
	if err != nil {
		fmt.Println("Error ", err);
		return make([]revil, 0)
	}
	defer rows.Close()

	revils := rowsToRevils(rows)
	return revils
}

func rowsToRevils(rows *sql.Rows) []revil {
>>>>>>> Improved looks of pages and changed functionality
	revils := make([]revil, 0)

	for rows.Next() {
		var url string
		var rtype string
		var comment string
		rows.Scan(&url, &rtype, &comment)
		revils = append(revils, revil{Type: rtype, Url: url, Comment: comment})
	}

	return revils
}

func getRevilInDatabase(row int) revil {
	rows, err := database.Query("select url, type, comment from revil LIMIT 1 OFFSET " + strconv.Itoa(row))
	if err != nil {
		fmt.Println("Error:", err)
		return *new(revil)
	}
	defer rows.Close()

	rows.Next()
	var url string
	var rtype string
	var comment string
	rows.Scan(&url, &rtype, &comment)
	return revil{Type: rtype, Url: url, Comment: comment}
}
