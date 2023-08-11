package database

func CreateUser(user UserInfo) (int, error) {
	sqlStmt, err := DB.Prepare(`INSERT INTO users(
		email,
		username,
		password,
		profile_picture)
	VALUES(?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}

	result, err := sqlStmt.Exec(user.Email, user.Username, user.Password, user.ProfilePicture)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func SelectAllUsers() ([]UserInfo, error) {
	rows, err := DB.Query("SELECT user_id, email, username, profile_picture FROM users")
	if err != nil {
		return nil, err
	}
	var users []UserInfo
	for rows.Next() {
		var user UserInfo
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.ProfilePicture)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func SelectUser(userID int) (*UserInfo, error) {
	row := DB.QueryRow("SELECT user_id, email, username, profile_picture FROM users WHERE user_id=?", userID)

	var user UserInfo
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.ProfilePicture)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserPassword(userID int) (*UserInfo, error) {
	row := DB.QueryRow("SELECT password FROM users WHERE user_id=?", userID)

	var user UserInfo
	err := row.Scan(&user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(userName string, email string, userID int) error {
	stmt, err := DB.Prepare(`
		UPDATE users SET
		email=?,
		username=?
		WHERE user_id=?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(email, userName, userID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(userID int) error {
	stmt, err := DB.Prepare(`
		DELETE FROM users
		WHERE user_id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}

func ReadUserLikedPosts(userID int) ([]int, error) {
	posts := make([]int, 0)

	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT PostId, CommentId
		FROM like WHERE Username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var like Like

		err = rows.Scan(&like.PostId, &like.CommentId)
		if err != nil {
			return nil, err
		}

		if like.PostId != 0 {
			post := like.PostId
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func ReadUserDislikedPosts(userID int) ([]int, error) {
	posts := make([]int, 0)

	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT PostId, CommentId
		FROM dislike WHERE Username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dislike Dislike

		err = rows.Scan(&dislike.PostId, &dislike.CommentId)
		if err != nil {
			return nil, err
		}

		if dislike.PostId != 0 {
			post := dislike.PostId
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func ReadUserCreatedPosts(userID int) ([]int, error) {
	posts := make([]int, 0)

	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT id FROM post WHERE Username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post

		err = rows.Scan(&post.Id)
		if err != nil {
			return nil, err
		}

		if post.Id != 0 {
			post := post.Id
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func GetUsername(userID int) (string, error) {
	var username string

	err := DB.QueryRow("SELECT username FROM users WHERE user_id = ?", userID).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func GetAvatar(username string) (string, error) {
	var avatar string

	err := DB.QueryRow("SELECT profile_picture FROM users WHERE username = ?", username).Scan(&avatar)
	if err != nil {
		return "", err
	}

	return avatar, nil
}
