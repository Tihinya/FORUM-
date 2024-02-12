package validation

import (
	"database/sql"
	"regexp"
)

func ValidateEmail(email string) bool {
	regex := `^[A-Za-z0-9~\x60!#$%^&*()_\-+={\[}\]|\\:;"'<,>.?/]{1,64}@[a-z]{1,255}\.[a-z]{1,63}$`
	match, _ := regexp.MatchString(regex, email)
	return len(email) > 0 && match
}

func GetUserID(db *sql.DB, email string, username string) (int, error) {
	var userID int
	var query string
	var args []interface{}

	if email != "" {
		query = "SELECT user_id FROM users WHERE email = ? OR username = ?"
		args = []interface{}{email, username}
	} else {
		query = "SELECT user_id FROM users WHERE username = ?"
		args = []interface{}{username}
	}

	err := db.QueryRow(query, args...).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User with the given email/username does not exist
			return 0, nil
		}
		// Other error occurred while querying the database
		return 0, err
	}

	// User with the given email/username exists, return their user ID
	return userID, nil
}

func ValidateUsername(username string) bool {

	bannedNames := []string{"admin", "Admin", "ADMIN", "HungryStepan"}

	for _, bannedName := range bannedNames {
		if username == bannedName {
			return false
		}
	}

	return true
}
func GetUserIdFromUserName(db *sql.DB, userName string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE username = ?", userName).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// User with the given email does not exist
			return 0, nil
		}
		// Other error occurred while querying the database
		return 0, err
	}
	// User with the given email exists, return their user ID
	return userID, nil
}
func ValidateRole(db *sql.DB, roleName string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM roles WHERE name = ?", roleName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func HasPendingRoleRequest(db *sql.DB, userId int) (bool, error) {
	// Query to check if the user has a pending role request
	query := `
		SELECT COUNT(*) FROM role_requests
		WHERE user_id = ?;
	`

	var count int
	err := db.QueryRow(query, userId).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func HasPendingPostReport(db *sql.DB, PostID int) (bool, error) {
	// Query to check if a post has a pending report
	query := `
		SELECT COUNT(*) FROM post_reports
		WHERE post_id = ? AND status = 'pending';
	`

	var count int
	err := db.QueryRow(query, PostID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetReportPostStatus(db *sql.DB, ReportID int) (string, error) {
	// Query to select the status of a report based on reportID
	query := "SELECT status FROM post_reports WHERE report_id = ?"

	var status string

	// Execute the query and scan the result into the 'status' variable
	err := db.QueryRow(query, ReportID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			// Report with the given reportID not found
			return "", nil
		}
		return "", err
	}

	return status, nil
}

func HasPost(db *sql.DB, PostID int) (bool, error) {
	// Query to check if the user has a pending role request
	query := `
		SELECT COUNT(*) FROM post
		WHERE id = ?;
	`

	var count int
	err := db.QueryRow(query, PostID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetUserName(db *sql.DB, userId int) (string, error) {
	// Query to retrieve the username for the given userId
	query := `
		SELECT username FROM users
		WHERE user_id = ?;
	`

	var username string
	err := db.QueryRow(query, userId).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func GetRoleName(db *sql.DB, RoleID int) (string, error) {
	// Query to retrieve the role name for the given RoleID
	query := `
		SELECT name FROM roles
		WHERE role_id = ?;
	`

	var roleName string
	err := db.QueryRow(query, RoleID).Scan(&roleName)
	if err != nil {
		return "", err
	}

	return roleName, nil
}
