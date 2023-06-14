package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var err error

// Run `go get github.com/mattn/go-sqlite3` in terminal to download db driver
func CreateTables() {
	db, err = sql.Open("sqlite3", "./database/database.db")
	checkErr(err)

	// Post database table
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS post (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
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
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			postId INTEGER NOT NULL,
			Content TEXT NOT NULL,
			Avatar TEXT,
			Username TEXT,
			CreationDate DATETIME,
			Likes INTEGER DEFAULT 0,
			Dislikes INTEGER DEFAULT 0,
			LastEdited DATETIME NULL
			FOREIGN KEY (postId) REFERENCES post(id)
		)
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	// Create another table:
	// ....

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
