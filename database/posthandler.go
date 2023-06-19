package database

import (
	"encoding/json"
	"fmt"
	"strconv"
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
		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)

		posts = append(posts, post)
	}

	// Convert posts to json
	jsonPosts, err := json.Marshal(posts)
	checkErr(err)

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
		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)

		posts = append(posts, post)
	}

	// Convert posts to json
	jsonPosts, err := json.Marshal(posts)
	checkErr(err)

	return jsonPosts
}

func UpdatePost(post Post, postID string) bool {
	stmt, _ := db.Prepare(`
		UPDATE post SET
			Title = ?,
			Content = ?,
			Categories = ?,
			LastEdited = ?
		WHERE id = ?
	`)

	if post.Id != 0 {
		postID = strconv.Itoa(post.Id)
	}

	// Checks if post with given ID exists in DB
	if !checkIfPostExist(postID) {
		return false
	}

	categoriesJSON, _ := json.Marshal(post.Categories)

	stmt.Exec(post.Title, post.Content, string(categoriesJSON), time.Now(), postID)

	return true
}

func DeletePost(postID string) bool {
	if !checkIfPostExist(postID) {
		return false
	}

	stmt, _ := db.Prepare(`
		DELETE FROM post WHERE ID = ?
	`)
	stmt.Exec(postID)

	return true
}

func checkIfPostExist(postID string) bool {
	err = db.QueryRow("SELECT 1 FROM post WHERE id = ?", postID).Scan(&postID)

	if err != nil {
		return false
	}

	return true
}
