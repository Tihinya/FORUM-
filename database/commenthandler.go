package database

import (
	"time"
)

func CreateCommentRow(comment Comment, postId int) (bool, error) {
	if !checkIfPostExist(postId) {
		return false, nil
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
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(postId, comment.Content, comment.UserInfo.ProfilePicture, comment.UserInfo.Username, comment.CreationDate)
	if err != nil {
		return false, err
	}

	return true, nil
}

func SelectComment(commentId int) ([]Comment, error) {
	var comments []Comment

	rows, err := DB.Query("SELECT * FROM comment where id = ?", commentId)
	if err != nil {
		return nil, err
	}

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
		if err != nil {
			return nil, err
		}

		comment.Likes, _ = getCommentLikes(comment.Id)
		comment.Dislikes, _ = getCommentDislikes(comment.Id)

		comment.UserInfo.ProfilePicture, _ = GetAvatar(comment.UserInfo.Username)

		comments = append(comments, comment)
	}

	return comments, nil
}

// GET all comments from comments table
func SelectAllComments(id int) ([]Comment, error) {
	var comments []Comment

	rows, err := DB.Query("SELECT * FROM comment where post_id = ?", id)
	if err != nil {
		return nil, err
	}

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
		if err != nil {
			return nil, err
		}

		comment.Likes, _ = getCommentLikes(comment.Id)
		comment.Dislikes, _ = getCommentDislikes(comment.Id)

		comment.UserInfo.ProfilePicture, _ = GetAvatar(comment.UserInfo.Username)

		comments = append(comments, comment)
	}

	return comments, nil
}

func UpdateComment(comment Comment, commentId int, username string) (bool, error) {
	if !checkCommentOwnership(commentId, username) {
		return false, nil
	}

	if !checkIfCommentExist(commentId) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		UPDATE comment SET
			content = ?,
			last_edited = ?
		WHERE id = ?
	`)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(comment.Content, time.Now(), commentId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteComment(commentID int, username string) (bool, error) {
	if !checkCommentOwnership(commentID, username) {
		return false, nil
	}

	if !checkIfCommentExist(commentID) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		DELETE FROM comment WHERE id = ?
	`)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(commentID)
	if err != nil {
		return false, err
	}

	err = deleteCommentLikes(commentID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func deletePostComments(postId int) error {
	stmt, err := DB.Prepare(`
		DELETE FROM comment WHERE post_id = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(postId)
	if err != nil {
		return err
	}

	return nil
}

func deleteCommentLikes(commentId int) error {
	stmt, err := DB.Prepare(`
		DELETE FROM like WHERE CommentId = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(commentId)
	if err != nil {
		return err
	}

	stmt, err = DB.Prepare(`
		DELETE FROM dislike WHERE CommentId = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(commentId)
	if err != nil {
		return err
	}

	return nil
}

func getPostCommentIds(postId int) ([]int, error) {
	var comments []int

	rows, err := DB.Query(`SELECT id FROM comment WHERE post_id = ?`, postId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment int

		err = rows.Scan(&comment)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// Checks if comment with given ID exists in DB
func checkIfCommentExist(commentId int) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM comment WHERE id=?)", commentId).Scan(&exists)

	return err == nil && exists
}

func checkCommentOwnership(commentId int, username string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM comment WHERE id = ? AND username = ?)", commentId, username).Scan(&exists)

	return err == nil && exists
}
