package database

import "database/sql"

func CreateRoleRequest(userId int, newRoleId int) error {
	// Prepare the INSERT statement
	stmt, err := DB.Prepare(`
		INSERT INTO role_requests (user_id, new_role_id)
		VALUES (?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, newRoleId)
	if err != nil {
		return err
	}

	return nil
}

func SelectAllPromotion() ([]RoleRequest, error) {
	var roleRequests []RoleRequest

	// Query to retrieve role requests
	query := `
		SELECT user_id, new_role_id
		FROM role_requests;
	`

	// Execute the query
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the results and append them to the roleRequests slice
	for rows.Next() {
		var roleRequest RoleRequest
		err := rows.Scan(&roleRequest.UserID, &roleRequest.NewRoleID)
		if err != nil {
			return nil, err
		}
		roleRequests = append(roleRequests, roleRequest)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roleRequests, nil
}

func DeleteRoleRequest(UserId int) error {
	query := `
		DELETE FROM role_requests
		WHERE user_id = ?;
	`

	_, err := DB.Exec(query, UserId)
	if err != nil {
		return err
	}

	return nil
}

func DemoteUser(RoleId int, UserId int) error {
	query := `
		UPDATE users SET role_id = ?
		WHERE user_id = ?;
	`

	_, err := DB.Exec(query, RoleId, UserId)
	if err != nil {
		return err
	}

	return nil
}

func GetNewRoleIDByUserID(UserID int) (int, error) {
	// Query to retrieve the new_role_id from the role_requests table
	query := `
		SELECT new_role_id FROM role_requests
		WHERE user_id = ?;
	`

	var newRoleID int
	err := DB.QueryRow(query, UserID).Scan(&newRoleID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No pending request found, return a default value or handle it accordingly
			return 0, nil
		}
		return 0, err
	}

	return newRoleID, nil
}
