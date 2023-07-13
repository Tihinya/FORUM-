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

func UndislikePost(postId int, username string) (bool, error) {
	if !checkIfPostExist(postId) {
		return false, nil
	}

	if !CheckIfPostDisliked(postId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		DELETE FROM dislike WHERE PostId = ? AND Username = ?
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

/*
func DislikeComment(commentId int, username string) (bool, error) {
	if !checkIfCommentExist(commentId) {
		return false, nil
	}

	if checkIfCommentDisliked(commentId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		INSERT INTO dislike (
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

func UndislikeComment(commentId int, username string) (bool, error) {
	if !checkIfCommentExist(commentId) {
		return false, nil
	}

	if !checkIfCommentDisliked(commentId, username) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		DELETE FROM dislike WHERE CommentId = ? AND Username = ?
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
*/

func getPostDislikes(postId int) (int, error) {
	var count int

	err := DB.QueryRow(`SELECT COUNT(*) FROM dislike WHERE PostId = ?`, postId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil

}

func CheckIfPostDisliked(postId int, username string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM dislike WHERE PostId = ? AND Username = ?)", postId, username).Scan(&exists)

	return err == nil && exists
}

func checkIfCommentDisliked(commentId int, username string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM dislike WHERE CommentId = ? AND Username = ?)", commentId, username).Scan(&exists)

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
