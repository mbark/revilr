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
		"CREATE TABLE revil (url TEXT NOT NULL, type TEXT, comment TEXT, date DATE DEFAULT (DATETIME('now','localtime')));",
		"CREATE TABLE user (username TEXT NOT NULL, password TEXT NOT NULL);",
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
	rows, err := database.Query("SELECT url, type, comment, date FROM revil")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		rowToRevil(rows).printRevil()
	}
}

func getAllRevilsInDatabase() []revil {
	rows, err := database.Query("SELECT url, type, comment, date FROM revil ORDER BY ROWID DESC")
	if err != nil {
		fmt.Println("Error:", err)
		return make([]revil, 0)
	}
	defer rows.Close()

	return rowsToRevils(rows)
}

func getRevilsOfType(rtype string) []revil {
	rows, err := database.Query("SELECT url, type, comment, date FROM revil WHERE type=? ORDER BY ROWID DESC", rtype)
	if err != nil {
		fmt.Println("Error ", err)
		return make([]revil, 0)
	}
	defer rows.Close()

	revils := rowsToRevils(rows)
	return revils
}

func rowsToRevils(rows *sql.Rows) []revil {
	revils := make([]revil, 0)

	for rows.Next() {
		revils = append(revils, rowToRevil(rows))
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
	return rowToRevil(rows)
}

func rowToRevil(row *sql.Rows) revil {
	var url string
	var rtype string
	var comment string
	var date string
	row.Scan(&url, &rtype, &comment, &date)
	return revil{Type: rtype, Url: url, Comment: comment, Date: date}
}

func getUser(username string) (user *User, err error) {
	rows, err := database.Query("select username, password from user WHERE username=?", username)
	if err != nil {
		return
	}
	defer rows.Close()

	rows.Next()
	user = rowToUser(rows)
	return
}

func rowToUser(row *sql.Rows) *User {
	var username string
	var password []byte
	row.Scan(&username, &password)
	return &User{Username: username, Password: password}
}

func createUser(user *User) error {
	stmt, err := database.Prepare("insert into user(username, password) values(?, ?)")
	if err != nil {
		return err
	}
	stmt.Exec()
	_, err = stmt.Exec(user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}
