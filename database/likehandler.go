package database

func LikePost(postId int, username string) (bool, error) {
	if !checkIfPostExist(postId) {
		return false, nil
	}

	if checkIfPostLiked(postId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		INSERT INTO like (
			PostId,
			Username
		) VALUES (?, ?)
	`)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(postId, username)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UnlikePost(postId int, username string) (bool, error) {
	if !checkIfPostExist(postId) {
		return false, nil
	}

	if !checkIfPostLiked(postId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		DELETE FROM like WHERE PostId = ? AND Username = ?
	`)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(postId, username)
	if err != nil {
		return false, err
	}

	return true, nil
}

func LikeComment(commentId int, username string) (bool, error) {
	if !checkIfCommentExist(commentId) {
		return false, nil
	}

	if checkIfCommentLiked(commentId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		INSERT INTO like (
			CommentId,
			Username
		) VALUES (?, ?)
	`)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(commentId, username)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UnlikeComment(commentId int, username string) (bool, error) {
	if !checkIfCommentExist(commentId) {
		return false, nil
	}

	if !checkIfCommentLiked(commentId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		DELETE FROM like WHERE CommentId = ? AND Username = ?
	`)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(commentId, username)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getPostLikes(postId int) (int, error) {
	var count int

	err := DB.QueryRow(`SELECT COUNT(*) FROM like WHERE PostId = ?`, postId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil

}

func checkIfPostLiked(postId int, username string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM like WHERE PostId = ? AND Username = ?)", postId, username).Scan(&exists)

	return err == nil && exists
}

func checkIfCommentLiked(commentId int, username string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM like WHERE CommentId = ? AND Username = ?)", commentId, username).Scan(&exists)

	return err == nil && exists
}

func Temp_selectLikes() ([]Like, error) {
	var likes []Like
	rows, err := DB.Query("SELECT * FROM like")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var like Like

		err = rows.Scan(&like.PostId, &like.CommentId, &like.Username)
		if err != nil {
			return nil, err
		}

		likes = append(likes, like)
	}

	return likes, nil
}
