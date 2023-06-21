package database

import (
	"time"
)

func CreateCommentRow(comment Comment, postId string) bool {

	if !checkIfPostExist(postId) {
		return false
	}

	stmt, err := db.Prepare(`
		INSERT INTO comment (
			PostId,
			Content,
			Avatar,
			Username,
			CreationDate
		) VALUES (?, ?, ?, ?, ?)
	`)
	checkErr(err)

	_, err = stmt.Exec(postId, comment.Content, comment.UserInfo.Avatar, comment.UserInfo.Username, comment.CreationDate)
	checkErr(err)

	return true
}

func SelectComment(commentId string) Comment {
	var comment Comment

	rows, err := db.Query("SELECT * FROM comment where Id = ?", commentId)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.Content,
			&comment.UserInfo.Avatar,
			&comment.UserInfo.Username,
			&comment.CreationDate,
			&comment.Likes,
			&comment.Dislikes,
			&comment.LastEdited,
		)
		checkErr(err)
	}

	return comment
}

// GET all comments from comments table
func SelectAllComments(id string) []Comment {
	var comments []Comment

	rows, err := db.Query("SELECT * FROM comment where PostId = ?", id)
	checkErr(err)

	for rows.Next() {
		var comment Comment

		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.Content,
			&comment.UserInfo.Avatar,
			&comment.UserInfo.Username,
			&comment.CreationDate,
			&comment.Likes,
			&comment.Dislikes,
			&comment.LastEdited,
		)
		checkErr(err)

		comments = append(comments, comment)
	}

	return comments
}

func UpdateComment(comment Comment, commentId string) bool {
	stmt, err := db.Prepare(`
		UPDATE comment SET
			Content = ?,
			LastEdited = ?
		WHERE id = ?
	`)
	checkErr(err)

	// Checks if comment with given ID exists in DB
	if !checkIfCommentExist(commentId) {
		return false
	}

	_, err = stmt.Exec(comment.Content, time.Now(), commentId)
	checkErr(err)

	return true
}

func DeleteComment(commentID string) bool {
	if !checkIfCommentExist(commentID) {
		return false
	}

	stmt, err := db.Prepare(`
		DELETE FROM comment WHERE ID = ?
	`)
	checkErr(err)

	_, err = stmt.Exec(commentID)
	checkErr(err)

	return true
}

func checkIfCommentExist(commentId string) bool {
	err = db.QueryRow("SELECT 1 FROM comment WHERE id=?", commentId).Scan(&commentId)

	if err != nil {
		return false
	}

	return true
}
