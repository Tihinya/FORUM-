package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/router"
	"forum/session"
	"log"
	"net/http"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	// Authentication here
	sessionToken, sessionTokenFound := checkForSessionToken(r)
	if !sessionTokenFound {
		returnMessageJSON(w, "Session token not found", http.StatusUnauthorized, "unauthorized")
		return
	}

	if !checkIfUserLoggedin(sessionToken) {
		returnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "unauthorized")
		return
	}

	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	username, err := database.GetUsername(userID)
	if err != nil {
		returnMessageJSON(w, "You are not logged in", http.StatusInternalServerError, "unauthorized")
		return
	}

	// Check if post is already disliked
	if database.CheckIfPostDisliked(postId, username) {
		returnMessageJSON(w, "Failed to like post. Post is already disliked.", http.StatusBadRequest, "error")
		return
	}

	existsLiked, err = database.LikePost(postId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		returnMessageJSON(w, "Liking post failed, post is already liked or does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Post successfully liked", http.StatusOK, "success")
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	// Authentication here
	sessionToken, sessionTokenFound := checkForSessionToken(r)
	if !sessionTokenFound {
		returnMessageJSON(w, "Session token not found", http.StatusUnauthorized, "unauthorized")
		return
	}

	if !checkIfUserLoggedin(sessionToken) {
		returnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "unauthorized")
		return
	}

	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	username, err := database.GetUsername(userID)
	if err != nil {
		returnMessageJSON(w, "You are not logged in", http.StatusInternalServerError, "unauthorized")
		return
	}

	existsLiked, err = database.UnlikePost(postId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		returnMessageJSON(w, "Unliking post failed, post is not liked or does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Post successfully unliked", http.StatusOK, "success")
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	// Authentication here
	sessionToken, sessionTokenFound := checkForSessionToken(r)
	if !sessionTokenFound {
		returnMessageJSON(w, "Session token not found", http.StatusUnauthorized, "unauthorized")
		return
	}

	if !checkIfUserLoggedin(sessionToken) {
		returnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "unauthorized")
		return
	}

	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	username, err := database.GetUsername(userID)
	if err != nil {
		returnMessageJSON(w, "You are not logged in", http.StatusInternalServerError, "unauthorized")
		return
	}

	// Check if commennt is already disliked
	if database.CheckIfCommentDisliked(commentId, username) {
		returnMessageJSON(w, "Failed to like comment. Comment is already disliked.", http.StatusBadRequest, "error")
		return
	}

	existsLiked, err = database.LikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		returnMessageJSON(w, "Liking comment failed, comment is already liked or does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully liked", http.StatusOK, "success")
}

func UnlikeComment(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	// Authentication here
	sessionToken, sessionTokenFound := checkForSessionToken(r)
	if !sessionTokenFound {
		returnMessageJSON(w, "Session token not found", http.StatusUnauthorized, "unauthorized")
		return
	}

	if !checkIfUserLoggedin(sessionToken) {
		returnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "unauthorized")
		return
	}

	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	username, err := database.GetUsername(userID)
	if err != nil {
		returnMessageJSON(w, "You are not logged in", http.StatusInternalServerError, "unauthorized")
		return
	}

	existsLiked, err = database.UnlikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		returnMessageJSON(w, "Unliking comment failed, comment is not liked or does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully unliked", http.StatusOK, "success")
}

func Temp_getLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	likes, err := database.Temp_selectLikes()
	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(likes)
}
