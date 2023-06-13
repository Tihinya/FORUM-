package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	UserInfo     UserInfo  `json:"userInfo"`
	CreationDate time.Time `json:"creationDate"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	Categories   []string  `json:"categories"`
}

type UserInfo struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
}

var db *sql.DB
var err error

// Run `go get github.com/mattn/go-sqlite3` in terminal to download db driver
func CreateTables() {
	db, err = sql.Open("sqlite3", "./database/database.db")
	checkErr(err)

	// Post database table
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			text TEXT NOT NULL,
			Avatar TEXT,
			Username TEXT,
			CreationDate DATETIME,
			Likes INTEGER DEFAULT 0,
			Dislikes INTEGER DEFAULT 0,
			Categories TEXT
		);
	`)
	checkErr(err)

	// Create another table:
	// ....

	_, err = stmt.Exec()
	checkErr(err)
}

func CreatePost(post Post) {
	categoriesJSON, err := json.Marshal(post.Categories)
	checkErr(err)
	fmt.Println(string(categoriesJSON))

	stmt, _ := db.Prepare(`
		INSERT INTO posts (
			title,
			text,
			Avatar,
			Username,
			CreationDate,
			Categories
		) VALUES (?, ?, ?, ?, ?, ?)
	`)
	stmt.Exec(post.Title, post.Text, post.UserInfo.Avatar, post.UserInfo.Username, post.CreationDate, string(categoriesJSON))

	fmt.Println("Post successfully created")
}

func SelectPost(id string) []byte {
	var posts []Post
	rows, err := db.Query("SELECT * FROM posts where id='" + (id) + "'")
	checkErr(err)

	for rows.Next() {
		var post Post

		rows.Scan(&post.Id, &post.Title, &post.Text, &post.UserInfo.Avatar, &post.UserInfo.Username, &post.CreationDate, &post.Likes, &post.Dislikes, &post.Categories)

		posts = append(posts, post)
	}

	// Convert posts to json
	jsonPosts, err := json.Marshal(posts)

	if err != nil {
		log.Fatal(err)
	}

	return jsonPosts
}

// GET all posts from posts table
func SelectAllPosts() []byte {
	var posts []Post
	rows, err := db.Query("SELECT * FROM posts")
	checkErr(err)

	for rows.Next() {
		var post Post
		rows.Scan(&post.Id, &post.Title, &post.Text, &post.UserInfo.Avatar, &post.UserInfo.Username, &post.CreationDate, &post.Likes, &post.Dislikes, &post.Categories)
		fmt.Println(post.Categories)
		posts = append(posts, post)
	}

	// Convert posts to json
	jsonPosts, err := json.Marshal(posts)

	if err != nil {
		log.Fatal(err, "ERRORORRORO")
	}
	return jsonPosts
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
