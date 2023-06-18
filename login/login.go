package login

import (
	"encoding/json"
	"forum/session"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type RegistrationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Simulating a database
var Users []User

func AddLogin(w http.ResponseWriter, userId int) {
	token := session.SessionStorage.CreateSession(userId)
	session.SessionStorage.SetCookie(token, w)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Validate inputs
	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Email and password are required",
		})
		return
	}
	// Check email format
	if !validateEmail(req.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Invalid email format",
		})
		return
	}
	// Check if email is already taken
	if emailTaken(req.Email) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Email or username already taken",
		})
		return
	}
	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Failed to encrypt password",
		})
		return
	}
	// Create a new user
	newUser := User{
		ID:       len(Users) + 1,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Save the new user to the database
	Users = append(Users, newUser)

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RegistrationResponse{
		Status:  "success",
		Message: "Registration successful",
	})
}

// Check if the email is already taken
func emailTaken(email string) bool {
	for _, user := range Users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func validateEmail(email string) bool {
	regex := `^[A-Za-z0-9~\x60!#$%^&*()_\-+={\[}\]|\\:;"'<,>.?/]{1,64}@[a-z]{1,255}\.[a-z]{1,63}$`
	match, _ := regexp.MatchString(regex, email)
	return len(email) > 0 && match
}

/*
curl -X GET -H "Content-Type: application/json" -d '{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/registration

------------------------------

curl -X GET -H "Content-Type: application/json" -d '{
  "username": "john_doe",z
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/registration

------------------------------

curl -X GET -H "Content-Type: application/json" -d '{
  "username": "john_doe",
  "email": "johnexample.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/registration

------------------------------

curl -X GET -H "Content-Type: application/json" -d '{
  "username": "jane_doe",
  "email": "john@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/registration

curl -X GET -H "Content-Type: application/json" -d '{
  "username": "john_doe",
  "email": "jane@example.com",
  "password": "password123",
  "password_confirmation": "password123"
}' -k https://localhost:8080/registration

*/
