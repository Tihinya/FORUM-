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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			profile_picture TEXT,
			username TEXT,
			creation_date DATETIME,
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			categories TEXT NULL,
			last_edited DATETIME NULL
		);
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = DB.Prepare(`
		CREATE TABLE IF NOT EXISTS category (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category TEXT
		);
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = DB.Prepare(`
		CREATE TABLE IF NOT EXISTS post_category (
			post_id INTEGER,
			category_id INTEGER,
			FOREIGN KEY (post_id) REFERENCES post(id),
			FOREIGN KEY (category_id) REFERENCES categories(id)
		);
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = DB.Prepare(`
		CREATE TABLE IF NOT EXISTS comment (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			profile_picture TEXT,
			username TEXT,
			creation_date DATETIME,
			likes INTEGER DEFAULT 0,
			dislikes INTEGER DEFAULT 0,
			last_edited DATETIME NULL,
			FOREIGN KEY (post_id) REFERENCES post(id)
		)
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = DB.Prepare(`
		CREATE TABLE IF NOT EXISTS like (
			post_id INTEGER DEFAULT 0, 
			comment_id INTEGER DEFAULT 0,
			username TEXT NOT NULL
		) 
	`) // PostId -> post_id
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = DB.Prepare(`
		CREATE TABLE IF NOT EXISTS dislike (
			post_id INTEGER DEFAULT 0,
			comment_id INTEGER DEFAULT 0,
			username TEXT NOT NULL
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
	checkErr(err)

	_, err = users.Exec()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
