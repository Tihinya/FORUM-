package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/login"
	"forum/router"
	"log"
	"net/http"
)

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
type UpdateUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
type DeleteUserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	login.Registration(w, r)
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

	var req UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UpdateUserResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	err = database.UpdateUser(req.Username, req.Email, userID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(UpdateUserResponse{
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

	json.NewEncoder(w).Encode(DeleteUserResponse{
		Status:  "success",
		Message: "User deleted successfully",
	})
}
