package database

import (
	"encoding/json"
	"fmt"
	"log"
)

func CreateCommentRow(comment Comment) {
	var postId int

	err = db.QueryRow("SELECT last_insert_rowid()").Scan(&postId)
	checkErr(err)

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

}

func SelectComment(id string) []byte {
	var comments []Comment
	rows, err := db.Query("SELECT * FROM comment where id='" + (id) + "'")
	checkErr(err)

	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.PostId, &comment.Id, &comment.Content, &comment.UserInfo.Avatar, &comment.UserInfo.Username, &comment.CreationDate, &comment.Likes, &comment.Dislikes, &comment.LastEdited)
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
func SelectAllComments() []byte {
	var comments []Comment
	rows, err := db.Query("SELECT * FROM comment")
	checkErr(err)

	for rows.Next() {
		var comment Comment
		rows.Scan(&comment.PostId, &comment.Id, &comment.Content, &comment.UserInfo.Avatar, &comment.UserInfo.Username, &comment.CreationDate, &comment.Likes, &comment.Dislikes, &comment.LastEdited)
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

func UpdateComment(comment Comment, commentID int) {

}

func DeleteComment(commentID int) {
}
