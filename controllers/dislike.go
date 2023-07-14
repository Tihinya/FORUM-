package controllers

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/router"
	"forum/session"
	"log"
	"net/http"
)

func DislikePost(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	// Check if post is already liked
	if database.CheckIfPostLiked(postId, username) {
		returnMessageJSON(w, "Failed to dislike post. Post is already liked.", http.StatusBadRequest, "error")
		return
	}

	existsLiked, err = database.DislikePost(postId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Disliking post failed, post with id %v is already disliked or does not exist", postId)
		return
	}

	returnMessageJSON(w, "Post successfully disliked", http.StatusOK, "success")
}

func UndislikePost(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	existsLiked, err = database.UndislikePost(postId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Undisliking post failed, post with id %v is not disliked or does not exist", postId)
		return
	}

	returnMessageJSON(w, "Post successfully undisliked", http.StatusOK, "success")
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	// Check if comment is already liked
	if database.CheckIfCommentLiked(commentId, username) {
		returnMessageJSON(w, "Failed to dislike comment. Comment is already liked.", http.StatusBadRequest, "error")
		return
	}

	existsLiked, err = database.DislikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Disliking comment failed, comment with id %v is already disliked or does not exist", commentId)
		return
	}

	returnMessageJSON(w, "Comment successfully disliked", http.StatusOK, "success")
}

func UndislikeComment(w http.ResponseWriter, r *http.Request) {
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
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	existsLiked, err = database.UndislikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Undisliking comment failed, comment with id %v is not disliked or does not exist", commentId)
		return
	}

	returnMessageJSON(w, "Comment successfully undisliked", http.StatusOK, "success")
}

func Temp_getDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	likes, err := database.Temp_selectDislikes()
	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(likes)
}
