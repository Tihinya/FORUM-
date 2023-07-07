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

	if register.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "User name cannot be empty ",
		})
		return
	}

	// Validate inputs
	if register.Email == "" || register.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Email and password are required",
		})
		return
	}
	// Check email format
	if !validation.ValidateEmail(register.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Invalid email format",
		})
		return
	}
	// Password check
	if register.Password != register.PasswordConfirmation {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Password confirmation does not match",
		})
		return
	}

	// Check if email is already taken
	exist, err := validation.GetUserID(database.DB, register.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exist != 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Email already taken",
		})
		return
	}

	id, err := login.AddUser(register.Username, register.Email, register.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//set session
	token := session.SessionStorage.CreateSession(id)
	session.SessionStorage.SetCookie(token, w)

	json.NewEncoder(w).Encode(database.Response{
		Status:  "success",
		Message: "Registration successful",
	})
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}
	if !validation.ValidateEmail(req.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Invalid email format",
		})
		return
	}

	// Check if email is already taken
	exist, err := validation.GetUserID(database.DB, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exist != 0 {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(database.Response{
			Status:  "error",
			Message: "Email already taken",
		})
		return
	}

	err = database.UpdateUser(req.Username, req.Email, userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(database.Response{
		Status:  "success",
		Message: "User updated successfully",
	})
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

	json.NewEncoder(w).Encode(database.Response{
		Status:  "success",
		Message: "User deleted successfully",
	})
}
