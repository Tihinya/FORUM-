package controllers

import (
	"encoding/json"
	"fmt"
	"forum/login"
	"forum/router"
	"net/http"
)

type GetUserResponse struct {
	Status string      `json:"status"`
	User   *login.User `json:"user,omitempty"`
}

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
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetUserResponse{
			Status: "error",
		})
		return
	}
	// Simulate retrieving the user from the database
	user := getUserFromDB(userID)
	fmt.Println(user)
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(GetUserResponse{
			Status: "error",
		})
		return
	}
	json.NewEncoder(w).Encode(GetUserResponse{
		Status: "success",
		User:   user,
	})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UpdateUserResponse{
			Status:  "error",
			Message: "Invalid user id",
		})
		return
	}
	// Simulate retrieving the user from the database
	user := getUserFromDB(userID)

	var req UpdateUserRequest
	error := json.NewDecoder(r.Body).Decode(&req)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UpdateUserResponse{
			Status:  "error",
			Message: "Invalid request body",
		})
		return
	}

	user.Username = req.Username
	user.Email = req.Email

	json.NewEncoder(w).Encode(UpdateUserResponse{
		Status:  "success",
		Message: "User updated successfully",
	})
}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(DeleteUserResponse{
			Status:  "error",
			Message: "User not found",
		})
		return
	}
	for i := 0; i < len(login.Users); i++ {
		if login.Users[i].ID == userID {
			login.Users = append(login.Users[:i], login.Users[i+1:]...)
		}
	}

	json.NewEncoder(w).Encode(DeleteUserResponse{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

func ReadUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(login.Users)

}
func getUserFromDB(id int) *login.User {
	for i := 0; i < len(login.Users); i++ {
		if login.Users[i].ID == id {
			return &login.Users[i]
		}
	}
	return nil
}
