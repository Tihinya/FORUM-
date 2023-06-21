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

	stmt, err := db.Prepare(`
		INSERT INTO post (
			Title,
			Content,
			Avatar,
			Username,
			CreationDate,
			Categories
		) VALUES (?, ?, ?, ?, ?, ?)
	`)
	checkErr(err)

	_, err = stmt.Exec(post.Title, post.Content, post.UserInfo.Avatar, post.UserInfo.Username, post.CreationDate, categoriesJSON)
	checkErr(err)
}

func SelectPost(id string) []Post {
	var categoriesString string
	var posts []Post

	rows, err := db.Query("SELECT * FROM post where id='" + (id) + "'")
	checkErr(err)

	for rows.Next() {
		var post Post

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.UserInfo.Avatar,
			&post.UserInfo.Username,
			&post.CreationDate,
			&post.Likes,
			&post.Dislikes,
			&categoriesString,
			&post.LastEdited,
		)
		checkErr(err)
		err = json.Unmarshal([]byte(categoriesString), &post.Categories)
		checkErr(err)
		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)
		posts = append(posts, post)
	}

	return posts
}

// GET all posts from posts table
func SelectAllPosts() []Post {
	var posts []Post
	rows, err := db.Query("SELECT * FROM post")
	checkErr(err)

	for rows.Next() {
		var post Post
		var categoriesString string

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.UserInfo.Avatar,
			&post.UserInfo.Username,
			&post.CreationDate,
			&post.Likes,
			&post.Dislikes,
			&categoriesString,
			&post.LastEdited,
		)
		checkErr(err)

		err = json.Unmarshal([]byte(categoriesString), &post.Categories)
		checkErr(err)
		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)

		posts = append(posts, post)
	}

	return posts
}

func UpdatePost(post Post, postID string) bool {
	if !checkIfPostExist(postID) {
		return false
	}

	stmt, err := db.Prepare(`
		UPDATE post SET
			Title = ?,
			Content = ?,
			Categories = ?,
			LastEdited = ?
		WHERE id = ?
	`)
	checkErr(err)

	if post.Id != 0 {
		postID = strconv.Itoa(post.Id)
	}

	categoriesJSON, err := json.Marshal(post.Categories)
	checkErr(err)

	_, err = stmt.Exec(post.Title, post.Content, string(categoriesJSON), time.Now(), postID)
	checkErr(err)

	return true
}

func DeletePost(postID string) bool {
	if !checkIfPostExist(postID) {
		return false
	}

	stmt, err := db.Prepare(`
		DELETE FROM post WHERE ID = ?
	`)
	checkErr(err)

	_, err = stmt.Exec(postID)
	checkErr(err)

	return true
}

func checkIfPostExist(postID string) bool {
	err = db.QueryRow("SELECT 1 FROM post WHERE id = ?", postID).Scan(&postID)

	if err != nil {
		return false
	}

	return true
}
