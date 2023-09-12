package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/login"
	"forum/router"
	"forum/session"
	"forum/validation"
	"log"
	"net/http"
	"strings"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var register database.UserInfo

	err := json.NewDecoder(r.Body).Decode(&register)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	// Trim username, email, and password
	register.Username = strings.TrimSpace(register.Username)
	register.Email = strings.TrimSpace(register.Email)
	register.Password = strings.TrimSpace(register.Password)
	register.PasswordConfirmation = strings.TrimSpace(register.PasswordConfirmation)

	if register.Username == "" {
		ReturnMessageJSON(w, "User name cannot be empty", http.StatusBadRequest, "error")
		return
	}
	// Validate inputs
	if register.Email == "" || register.Password == "" {
		ReturnMessageJSON(w, "Email and password are required", http.StatusBadRequest, "error")
		return
	}
	// Check email format
	if !validation.ValidateEmail(register.Email) {
		ReturnMessageJSON(w, "Invalid email format", http.StatusBadRequest, "error")
		return
	}
	// Password check
	if register.Password != register.PasswordConfirmation {
		ReturnMessageJSON(w, "Password confirmation does not match", http.StatusBadRequest, "error")
		return
	}

	// Check if email is already taken
	exist, err := validation.GetUserIdFromEmail(database.DB, register.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exist != 0 {
		ReturnMessageJSON(w, "Email already taken", http.StatusConflict, "error")
		return
	}

	exist, err = validation.GetUserIdFromUserName(database.DB, register.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exist != 0 {
		ReturnMessageJSON(w, "Username already taken", http.StatusConflict, "error")
		return
	}

	id, err := login.AddUser(register.Username, register.Email, register.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set session
	token := session.SessionStorage.CreateSession(id)
	session.SessionStorage.SetCookie(token, w)

	ReturnMessageJSON(w, "Registration successful", http.StatusOK, "success")
}

func ReadUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := database.SelectUser(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func ReadUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := database.SelectAllUsers()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var req database.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	if !validation.ValidateEmail(req.Email) {
		ReturnMessageJSON(w, "Invalid email format", http.StatusBadRequest, "error")
		return
	}

	// Check if email is already taken
	exist, err := validation.GetUserIdFromEmail(database.DB, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exist != 0 {
		ReturnMessageJSON(w, "Email already taken", http.StatusConflict, "error")
		return
	}

	err = database.UpdateUser(req.Username, req.Email, userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ReturnMessageJSON(w, "User updated successfully", http.StatusOK, "success")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = database.DeleteUser(userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ReturnMessageJSON(w, "User deleted successfully", http.StatusOK, "success")
}

func ReadUserLikedPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	UserId := getUserId(r)

	posts, err := database.ReadUserLikedPosts(UserId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func ReadUserDislikedPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	UserId := getUserId(r)

	posts, err := database.ReadUserDislikedPosts(UserId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func ReadUserLikedComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	UserId := getUserId(r)

	comments, err := database.ReadUserLikedComments(UserId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func ReadUserDislikedComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	UserId := getUserId(r)

	comments, err := database.ReadUserDislikedComments(UserId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func ReadUserCreatedPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	UserId := getUserId(r)

	posts, err := database.ReadUserCreatedPosts(UserId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(posts)
}
