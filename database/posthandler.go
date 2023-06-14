package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func CreatePost(post Post) {
	categoriesJSON, err := json.Marshal(post.Categories)
	checkErr(err)

	stmt, _ := db.Prepare(`
		INSERT INTO post (
			Title,
			Content,
			Avatar,
			Username,
			CreationDate,
			Categories
		) VALUES (?, ?, ?, ?, ?, ?)
	`)
	stmt.Exec(post.Title, post.Content, post.UserInfo.Avatar, post.UserInfo.Username, post.CreationDate, categoriesJSON)

	fmt.Println("Post successfully created")
}

func SelectPost(id string) []byte {
	var posts []Post
	rows, err := db.Query("SELECT * FROM post where id='" + (id) + "'")
	checkErr(err)

	for rows.Next() {
		var post Post
		var categoriesString string
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserInfo.Avatar, &post.UserInfo.Username, &post.CreationDate, &post.Likes, &post.Dislikes, &categoriesString, &post.LastEdited)
		err = json.Unmarshal([]byte(categoriesString), &post.Categories)
		checkErr(err)
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
	rows, err := db.Query("SELECT * FROM post")
	checkErr(err)

	for rows.Next() {
		var post Post
		var categoriesString string
		rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserInfo.Avatar, &post.UserInfo.Username, &post.CreationDate, &post.Likes, &post.Dislikes, &categoriesString, &post.LastEdited)
		err = json.Unmarshal([]byte(categoriesString), &post.Categories)
		checkErr(err)
		posts = append(posts, post)
	}

	// Convert posts to json
	jsonPosts, err := json.Marshal(posts)

	if err != nil {
		log.Fatal(err, "ERRORORRORO")
	}
	return jsonPosts
}

func UpdatePost(post Post, postID int) {
	stmt, _ := db.Prepare(`
		UPDATE post SET
			Title = ?,
			Text = ?,
			Tategories = ?,
			LastEdited = ?
		WHERE id = ?
	`)

	if post.Id != 0 {
		postID = post.Id
	}

	categoriesJSON, _ := json.Marshal(post.Categories)

	stmt.Exec(post.Title, post.Content, string(categoriesJSON), time.Now(), postID)

	fmt.Println("Post successfully updated")
}

func DeletePost(postID int) {
	stmt, _ := db.Prepare(`
		DELETE FROM post WHERE ID = ?
	`)

	stmt.Exec(postID)
}
