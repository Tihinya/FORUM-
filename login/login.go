package login

import (
	"database/sql"
	"encoding/json"
	"forum/config"
	"forum/database"
	"forum/session"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type RegistrationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func AddLogin(w http.ResponseWriter, userId int) {
	token := session.SessionStorage.CreateSession(userId)
	session.SessionStorage.SetCookie(token, w)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	var register database.UserInfo

	err := json.NewDecoder(r.Body).Decode(&register)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Validate inputs
	if register.Email == "" || register.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Email and password are required",
		})
		return
	}
	// Check email format
	if !validateEmail(register.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Invalid email format",
		})
		return
	}
	// Password check
	if register.Password != register.PasswordConfirmation {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Password confirmation does not match",
		})
		return
	}

	// Check if email is already taken
	exist, err := checkEmailExists(database.DB, register.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exist {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Email or username already taken",
		})
		return
	}
	user := database.UserInfo{
		ProfilePicture: config.Config.ProfilePicture,
		Username:       register.Username,
		Email:          register.Email,
		Password:       register.Password,
	}

	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	id, err := database.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.ID = id

	//set session
	token := session.SessionStorage.CreateSession(id)
	session.SessionStorage.SetCookie(token, w)

	json.NewEncoder(w).Encode(RegistrationResponse{
		Status:  "success",
		Message: "Registration successful",
	})
}

func validateEmail(email string) bool {
	regex := `^[A-Za-z0-9~\x60!#$%^&*()_\-+={\[}\]|\\:;"'<,>.?/]{1,64}@[a-z]{1,255}\.[a-z]{1,63}$`
	match, _ := regexp.MatchString(regex, email)
	return len(email) > 0 && match
}
func checkEmailExists(db *sql.DB, email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
