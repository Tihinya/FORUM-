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