package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB  *sql.DB
	err error
)

// Open database and assign global DB the local DB value
func OpenDatabase(dataSourceName string) error {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}
	DB = db // Assign global DB the local DB value

	CreateTables()

	return nil
}

func CreateTables() {
	// Post database table
	stmt, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS post (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			image BLOB,
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

	// Create roles table
	roles, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS roles (
			role_id  INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		);
	`)
	checkErr(err)

	_, err = roles.Exec()
	checkErr(err)

	// Create users table
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
			gender TEXT,
			age TEXT,
			profile_picture BLOB,
			role_id INTEGER,
			FOREIGN KEY (role_id) REFERENCES roles(role_id)
		);
	`)
	checkErr(err)

	_, err = users.Exec()
	checkErr(err)

	// Create the role_requests table
	roleRequestsStmt, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS role_requests (
			request_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			new_role_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (new_role_id) REFERENCES roles(role_id)
			);
	`)
	checkErr(err)
	_, err = roleRequestsStmt.Exec()
	checkErr(err)

	stmt, err = DB.Prepare(`
		CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			parent_object_id INTEGER NOT NULL,
			related_object_type TEXT NOT NULL,
			related_object_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			status TEXT DEFAULT "unread" NOT NULL,
			creation_date DATETIME
		);
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)

	// Create the post_report table
	postReportStmt, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS post_reports (
			report_id INTEGER PRIMARY KEY AUTOINCREMENT,
			message TEXT,
			response TEXT,
			status TEXT CHECK (status IN ('pending', 'approved', 'rejected')),
			seen BOOL,
			user_id INTEGER,
			post_id INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(user_id),
			FOREIGN KEY (post_id) REFERENCES post(id)
		);
	`)
	checkErr(err)
	_, err = postReportStmt.Exec()
	checkErr(err)

	stmt, err = DB.Prepare(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			message_id INTEGER PRIMARY KEY AUTOINCREMENT,
			message TEXT NOT NULL,
			created_at DATETIME,
			receiver_id INTEGER,
			sender_id INTEGER,
			FOREIGN KEY (receiver_id) REFERENCES users(user_id),
			FOREIGN KEY (sender_id) REFERENCES users(user_id)
		);
	`)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
