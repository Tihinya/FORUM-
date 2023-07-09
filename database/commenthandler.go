package database

import (
	"time"
)

func CreateCommentRow(comment Comment, postId int) bool {

	if !checkIfPostExist(postId) {
		return false
	}

	stmt, err := DB.Prepare(`
		INSERT INTO comment (
			post_id,
			content,
			profile_picture,
			username,
			creation_date
		) VALUES (?, ?, ?, ?, ?)
	`)
	checkErr(err)

	_, err = stmt.Exec(postId, comment.Content, comment.UserInfo.ProfilePicture, comment.UserInfo.Username, comment.CreationDate)
	checkErr(err)

	return true
}

func SelectComment(commentId string) Comment {
	var comment Comment

	rows, err := DB.Query("SELECT * FROM comment where id = ?", commentId)
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.Content,
			&comment.UserInfo.ProfilePicture,
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

	rows, err := DB.Query("SELECT * FROM comment where post_id = ?", id)
	checkErr(err)

	for rows.Next() {
		var comment Comment

		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.Content,
			&comment.UserInfo.ProfilePicture,
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

func UpdateComment(comment Comment, commentId int) bool {
	stmt, err := DB.Prepare(`
		UPDATE comment SET
			content = ?,
			last_edited = ?
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

func DeleteComment(commentID int) bool {
	if !checkIfCommentExist(commentID) {
		return false
	}

	stmt, err := DB.Prepare(`
		DELETE FROM comment WHERE id = ?
	`)
	checkErr(err)

	_, err = stmt.Exec(commentID)
	checkErr(err)

	return true
}

func deletePostComments(postId int) {
	stmt, err := DB.Prepare(`
		DELETE FROM comment WHERE post_id = ?
	`)
	checkErr(err)

	_, err = stmt.Exec(postId)
	checkErr(err)
}

func checkIfCommentExist(commentId int) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM comment WHERE id=?)", commentId).Scan(&exists)

	return err == nil && exists
}
