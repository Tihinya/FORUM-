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
	// Regex pattern for username validation
	// Allows alphanumeric characters (a-z, A-Z, 0-9) and underscores
	// Must start with a letter
	// Must be between 3 and 16 characters in length
	pattern := "^[a-zA-Z][a-zA-Z0-9_]{2,15}$"

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Check if the username matches the pattern
	return regex.MatchString(username)
}
