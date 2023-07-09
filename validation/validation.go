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

func GetUserID(db *sql.DB, email string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM users WHERE email = ?", email).Scan(&userID)
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
