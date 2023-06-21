package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var err error

// Run `go get github.com/mattn/go-sqlite3` in terminal to download db driver
func CreateTables() {
	DB, err = sql.Open("sqlite3", "./database/database.db")
	checkErr(err)

	// Post database table
	stmt, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS post (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Title TEXT NOT NULL,
			Content TEXT NOT NULL,
			Avatar TEXT,
			Username TEXT,
			CreationDate DATETIME,
			Likes INTEGER DEFAULT 0,
			Dislikes INTEGER DEFAULT 0,
			Categories TEXT,
			LastEdited DATETIME NULL
		);
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	// Create another table:
	users, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			ID INTEGER PRIMARY KEY AUTOINCREMENT,
			Email TEXT,
			Username TEXT,
			Password TEXT,
			Avatar TEXT
		);
	`)
	checkErr(err)
	_, err = users.Exec()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
