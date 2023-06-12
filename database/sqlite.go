package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Run `go get github.com/mattn/go-sqlite3` in terminal to download db driver
func CreateTables() {
	db, err := sql.Open("sqlite3", "./database/database.db")
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			text TEXT,
			creationDate DATETIME
		);
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
