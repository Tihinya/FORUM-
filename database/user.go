package database

import "database/sql"

func CreateUser(user UserInfo) (int, error) {
	sqlStmt, err := DB.Prepare(`INSERT INTO users (
		email,
		username,
		password,
		profile_picture,
		role_id)
	VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}

	result, err := sqlStmt.Exec(user.Email, user.Username, user.Password, user.ProfilePicture, user.RoleID)
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

func UpdateUser(userName string, email string, roleID int, userID int) error {
	stmt, err := DB.Prepare(`
		UPDATE users SET
		email=?,
		username=?,
		role_id=?
		WHERE user_id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(email, userName, roleID, userID)
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

func GenerateDefaultRoles() error {
	// Check if the roles already exist in the database
	row := DB.QueryRow("SELECT COUNT(*) FROM roles")
	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	// If roles already exist, return without creating them again
	if count > 0 {
		return nil
	}

	// Insert the default roles into the roles table
	roles := []string{"user", "moderator", "admin"}
	stmt, err := DB.Prepare("INSERT INTO roles (name) VALUES (?)")
	if err != nil {
		return err
	}

	for _, role := range roles {
		_, err = stmt.Exec(role)
		if err != nil {
			return err
		}
	}

	return nil
}
func GetRoleId(roleName string) (int, error) {
	row := DB.QueryRow("SELECT role_id FROM roles WHERE name = ?", roleName)

	var roleId int
	err := row.Scan(&roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Role with the specified name not found
			return 0, nil
		}
		return 0, err
	}

	return roleId, nil
}
