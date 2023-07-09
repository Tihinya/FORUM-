package database

import (
	"encoding/json"
)

func LikePost(postId int, username string) bool {
	if !checkIfPostExist(postId) {
		return false
	}

	if !checkIfPostLiked(postId, username) {
		return false
	}

	stmt, err := DB.Prepare(`
		INSERT INTO like (
			PostId,
			Username
		) VALUES (?, ?)
	`)
	checkErr(err)

	stmt.Exec(postId, username)
	checkErr(err)

	return true
}

func UnlikePost(postId int, username string) bool {
	if !checkIfPostExist(postId) {
		return false
	}

	if !checkIfPostLiked(postId, username) {
		return false
	}

	stmt, _ := DB.Prepare(`
		DELETE FROM like WHERE PostId = ? AND Username = ?
	`)
	stmt.Exec(postId, username)

	return true
}

func LikeComment(commentId int, username string) bool {
	if !checkIfCommentExist(commentId) {
		return false
	}

	if !checkIfCommentLiked(commentId, username) {
		return false
	}

	stmt, _ := DB.Prepare(`
		INSERT INTO like (
			CommentId,
			Username
		) VALUES (?, ?)
	`)
	stmt.Exec(commentId, username)

	return true
}

func UnlikeComment(commentId int, username string) bool {
	if !checkIfCommentExist(commentId) {
		return false
	}

	if !checkIfCommentLiked(commentId, username) {
		return false
	}

	stmt, _ := DB.Prepare(`
		DELETE FROM like WHERE CommentId = ? AND Username = ?
	`)
	stmt.Exec(commentId, username)

	return true
}

func checkIfPostLiked(postId int, username string) bool {
	err = DB.QueryRow("SELECT 1 FROM like WHERE PostId = ? AND Username = ?", postId, username).Scan(&username)

	if err == nil {
		return false
	}

	return true

}

func checkIfCommentLiked(commentId int, username string) bool {
	err = DB.QueryRow("SELECT 1 FROM like WHERE CommentId = ? AND Username = ?", commentId, username).Scan(&username)

	if err == nil {
		return false
	}

	return true

}

func Temp_selectLikes() []byte {
	var likes []Like
	rows, err := DB.Query("SELECT * FROM like")
	checkErr(err)

	for rows.Next() {
		var like Like

		rows.Scan(&like.Id, &like.PostId, &like.CommentId, &like.Username)

		likes = append(likes, like)
	}

	// Convert posts to json
	jsonLikes, err := json.Marshal(likes)
	checkErr(err)

	return jsonLikes
}
