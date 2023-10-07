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
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
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
	exist, err := validation.GetUserID(database.DB, register.Email, "")
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
	roleId, err := database.GetRoleId("user")
	if err != nil {
		log.Println(err)
	}
	id, err := login.AddUser(register.Username, register.Email, register.Password, roleId)
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
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}
	json.NewEncoder(w).Encode(users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req database.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	//Get user ID
	userID, err := router.GetFieldInt(r, "id")
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}

	// Get user data
	user, err := database.SelectUser(userID)
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}
	// Get role ID
	roleId := user.RoleID
	if req.Role != "" {
		roleId, err = database.GetRoleId(req.Role)
		if err != nil {
			ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		}
	}
	// Check email
	if req.Email != "" {
		if !validation.ValidateEmail(req.Email) {
			ReturnMessageJSON(w, "Invalid email format", http.StatusBadRequest, "error")
			return
		}
		exist, err := validation.GetUserID(database.DB, req.Email, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if exist != 0 {
			ReturnMessageJSON(w, "Email already taken", http.StatusConflict, "error")
			return
		}
	}
	// Check username
	if req.Username != "" {
		if !validation.ValidateUsername(req.Username) {
			ReturnMessageJSON(w, "This name cannot be used", http.StatusBadRequest, "error")
			return
		}
	}
	// Check role
	if req.Role != "" {
		if roleId == user.RoleID {
			ReturnMessageJSON(w, "Role cannot be changed", http.StatusBadRequest, "error")
			return
		}
		// Check if role exist
		exist, err := validation.ValidateRole(database.DB, req.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !exist {
			ReturnMessageJSON(w, "Role doesn't exist ", http.StatusBadRequest, "error")
			return
		}

	}

	err = database.UpdateUser(req.Username, req.Email, roleId, userID)
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "User updated", http.StatusOK, "success")
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

func ReadUserCommentdPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	UserId := getUserId(r)
	posts, err := database.ReadUserCommentsPosts(UserId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(posts)

}
