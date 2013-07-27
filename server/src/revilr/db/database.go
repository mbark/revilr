package db

import (
	"revilr/user"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var database *sql.DB

func OpenConnection() (db *sql.DB, err error) {
	dbPath := "./revils.db"

	if exists, _ := fileExists(dbPath); exists {
		db, err = sql.Open("sqlite3", dbPath)
	} else {
		db, err = instantiateDatabase(dbPath)
	}
	database = db

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
		"CREATE TABLE revil (url TEXT NOT NULL, type TEXT, comment TEXT, user TEXT NOT NULL, date DATE DEFAULT (DATETIME('now','localtime')));",
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

func InsertIntoDatabase(rev user.Revil, usr user.User) error {
	stmt, err := database.Prepare("insert into revil(url, type, comment, user) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	stmt.Exec()
	_, err = stmt.Exec(rev.Url, rev.Type, rev.Comment, usr.Username)
	return err
}

func GetAllRevilsInDatabase(usr user.User) []user.Revil {
	rows, err := database.Query("SELECT url, type, comment, date FROM revil WHERE user = ? ORDER BY ROWID DESC", usr.Username)
	if err != nil {
		fmt.Println(err)
		return make([]user.Revil, 0)
	}
	defer rows.Close()

	return rowsToRevils(rows)
}

func GetRevilsOfType(rtype string, usr user.User) []user.Revil {
	rows, err := database.Query("SELECT url, type, comment, date FROM revil WHERE type=? AND user=? ORDER BY ROWID DESC", rtype, usr.Username)
	if err != nil {
		fmt.Println("Error ", err)
		return make([]user.Revil, 0)
	}
	defer rows.Close()

	revils := rowsToRevils(rows)
	return revils
}

func rowsToRevils(rows *sql.Rows) []user.Revil {
	revils := make([]user.Revil, 0)

	for rows.Next() {
		revils = append(revils, rowToRevil(rows))
	}

	return revils
}

func rowToRevil(row *sql.Rows) user.Revil {
	var url string
	var rtype string
	var comment string
	var date string
	row.Scan(&url, &rtype, &comment, &date)
	return user.Revil{Type: rtype, Url: url, Comment: comment, Date: date}
}

func FindUser(username string) (user *user.User, err error) {
	rows, err := database.Query("select username, password from user WHERE username=?", username)
	if err != nil {
		return
	}
	defer rows.Close()

	rows.Next()
	user = rowToUser(rows)
	return
}

func rowToUser(row *sql.Rows) *user.User {
	var username string
	var password []byte
	row.Scan(&username, &password)
	return user.NewType(username, password)
}

func CreateUser(user *user.User) error {
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
