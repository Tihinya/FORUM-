package database

func DislikePost(postId int, username string) (bool, error) {
	if !checkIfPostExist(postId) {
		return false, nil
	}

	if CheckIfPostDisliked(postId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		INSERT INTO dislike (
			post_id,
			username
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

func UndislikePost(postId int, username string) (bool, error) {
	if !checkIfPostExist(postId) {
		return false, nil
	}

	if !CheckIfPostDisliked(postId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		DELETE FROM dislike WHERE post_id = ? AND username = ?
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

func DislikeComment(commentId int, username string) (bool, error) {
	if !checkIfCommentExist(commentId) {
		return false, nil
	}

	if CheckIfCommentDisliked(commentId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		INSERT INTO dislike (
			comment_id,
			username
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

func UndislikeComment(commentId int, username string) (bool, error) {
	if !checkIfCommentExist(commentId) {
		return false, nil
	}

	if !CheckIfCommentDisliked(commentId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		DELETE FROM dislike WHERE comment_id = ? AND username = ?
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

func getPostDislikes(postId int) (int, error) {
	var count int

	err := DB.QueryRow(`SELECT COUNT(*) FROM dislike WHERE post_id = ?`, postId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getCommentDislikes(commentId int) (int, error) {
	var count int

	err := DB.QueryRow(`SELECT COUNT(*) FROM dislike WHERE comment_id = ?`, commentId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CheckIfPostDisliked(postId int, username string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM dislike WHERE post_id = ? AND username = ?)", postId, username).Scan(&exists)

	return err == nil && exists
}

func CheckIfCommentDisliked(commentId int, username string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM dislike WHERE comment_id = ? AND username = ?)", commentId, username).Scan(&exists)

	return err == nil && exists
}

func Temp_selectDislikes() ([]Dislike, error) {
	var Dislikes []Dislike
	rows, err := DB.Query("SELECT * FROM dislike")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var Dislike Dislike

		err = rows.Scan(&Dislike.PostId, &Dislike.CommentId, &Dislike.Username)
		if err != nil {
			return nil, err
		}

		Dislikes = append(Dislikes, Dislike)
	}

	return Dislikes, nil
}
