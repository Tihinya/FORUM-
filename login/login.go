package login

import (
	"encoding/json"
	"forum/session"
	"net/http"
)

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

	// Check if email or username is already taken (perform your checks here)
	emailTaken := false
	usernameTaken := false
	if emailTaken || usernameTaken {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(RegistrationResponse{
			Status:  "error",
			Message: "Email or username already taken",
		})
		return
	}

	// Perform registration logic (create user in the database, etc.)
	// ...

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RegistrationResponse{
		Status:  "success",
		Message: "Registration successful",
	})
}
