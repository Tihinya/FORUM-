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

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS comment (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			PostId INTEGER NOT NULL,
			Content TEXT NOT NULL,
			Avatar TEXT,
			Username TEXT,
			CreationDate DATETIME,
			Likes INTEGER DEFAULT 0,
			Dislikes INTEGER DEFAULT 0,
			LastEdited DATETIME NULL,
			FOREIGN KEY (PostId) REFERENCES post(Id)
		)
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS like (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			PostId INTEGER DEFAULT 0,
			CommentId INTEGER DEFAULT 0,
			Username TEXT NOT NULL
		)
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	// Create another table:
	users, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			user_id  INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT,
			username TEXT,
			password TEXT,
			profile_picture BLOB
		);
	`)
	if err != nil {
		log.Println(err)
	}

	_, err = users.Exec()
	if err != nil {
		log.Println(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
