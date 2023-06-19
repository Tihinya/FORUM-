package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func CreateCommentRow(comment Comment, postId string) bool {

	if !checkIfPostExist(postId) {
		return false
	}

	stmt, _ := db.Prepare(`
		INSERT INTO comment (
			PostId,
			Content,
			Avatar,
			Username,
			CreationDate
		) VALUES (?, ?, ?, ?, ?)
	`)

	stmt.Exec(postId, comment.Content, comment.UserInfo.Avatar, comment.UserInfo.Username, comment.CreationDate)

	fmt.Println("Comment successfully created")
	return true

}

func SelectComment(commentId string) []byte {
	var comments []Comment

	rows, err := db.Query("SELECT * FROM comment where Id = ?", commentId)
	checkErr(err)

	for rows.Next() {
		var comment Comment

		rows.Scan(&comment.Id, &comment.PostId, &comment.Content, &comment.UserInfo.Avatar, &comment.UserInfo.Username, &comment.CreationDate, &comment.Likes, &comment.Dislikes, &comment.LastEdited)
		checkErr(err)

		comments = append(comments, comment)
	}

	// Convert comments to json
	jsoncomments, err := json.Marshal(comments)

	if err != nil {
		log.Fatal(err)
	}

	return jsoncomments
}

// GET all comments from comments table
func SelectAllComments(id string) []byte {
	var comments []Comment

	rows, err := db.Query("SELECT * FROM comment where PostId = ?", id)
	checkErr(err)

	for rows.Next() {
		var comment Comment

		rows.Scan(&comment.Id, &comment.PostId, &comment.Content, &comment.UserInfo.Avatar, &comment.UserInfo.Username, &comment.CreationDate, &comment.Likes, &comment.Dislikes, &comment.LastEdited)
		checkErr(err)

		comments = append(comments, comment)
	}

	// Convert comments to json
	jsoncomments, err := json.Marshal(comments)

	if err != nil {
		log.Fatal(err)
	}

	return jsoncomments
}

func UpdateComment(comment Comment, commentId string) bool {
	stmt, _ := db.Prepare(`
		UPDATE comment SET
			Content = ?,
			LastEdited = ?
		WHERE id = ?
	`)

	// Checks if comment with given ID exists in DB
	if !checkIfCommentExist(commentId) {
		return false
	}

	stmt.Exec(comment.Content, time.Now(), commentId)

	return true
}

func DeleteComment(commentID string) bool {
	if !checkIfCommentExist(commentID) {
		return false
	}

	stmt, _ := db.Prepare(`
		DELETE FROM comment WHERE ID = ?
	`)
	stmt.Exec(commentID)

	return true
}

func checkIfCommentExist(commentId string) bool {
	err = db.QueryRow("SELECT 1 FROM comment WHERE id=?", commentId).Scan(&commentId)

	if err != nil {
		return false
	}

	return true
}
