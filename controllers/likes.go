package controllers

import (
	"fmt"
	"forum/database"
	"forum/router"
	"net/http"
)

var username = "brozkie"

func LikePost(w http.ResponseWriter, r *http.Request) {
	postId, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	// Authentication here maybe

	if !database.LikePost(postId, username) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Liking post failed, post with id %v is already liked or does not exist", postId)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post liked")
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	postId, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	// Authentication here maybe

	database.UnlikePost(postId, username)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post unliked")
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	// Authentication here maybe

	if !database.LikeComment(commentId, username) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Liking comment failed, comment with id %v is already liked or does not exist", commentId)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Comment liked")
}

func UnlikeComment(w http.ResponseWriter, r *http.Request) {

}

func Temp_getLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	likes := database.Temp_selectLikes()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", likes)
}
