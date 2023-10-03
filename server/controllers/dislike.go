package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/router"
	"log"
	"net/http"
)

func DislikePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var existsLiked bool

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	username := getUsername(r)

	// Check if post is already liked
	if database.CheckIfPostLiked(postId, username) {
		ReturnMessageJSON(w, "Failed to dislike post. Post is already liked.", http.StatusBadRequest, "error")
		return
	}

	existsLiked, err = database.DislikePost(postId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		ReturnMessageJSON(w, "Disliking post failed, post is already disliked or does not exist", http.StatusBadRequest, "error")
		return
	}

	err = database.CreateNotification(postId, "post", postId, "dislike")
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "Post successfully disliked", http.StatusOK, "success")
}

func UndislikePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var existsLiked bool

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	username := getUsername(r)
	existsLiked, err = database.UndislikePost(postId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		ReturnMessageJSON(w, "Undisliking post failed, post is not disliked or does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Post successfully undisliked", http.StatusOK, "success")
}

func DislikeComment(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	username := getUsername(r)

	// Check if comment is already liked
	if database.CheckIfCommentLiked(commentId, username) {
		ReturnMessageJSON(w, "Failed to dislike comment. Comment is already liked.", http.StatusBadRequest, "error")
		return
	}

	existsLiked, err = database.DislikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		ReturnMessageJSON(w, "Disliking comment failed, comment is already disliked or does not exist", http.StatusBadRequest, "error")
		return
	}

	postId, err := database.GetPostIdFromCommentId(commentId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal error", http.StatusInternalServerError, "error")
		return
	}

	err = database.CreateNotification(postId, "comment", commentId, "dislike")
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "Comment successfully disliked", http.StatusOK, "success")
}

func UndislikeComment(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	username := getUsername(r)
	existsLiked, err = database.UndislikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		ReturnMessageJSON(w, "Undisliking comment failed, comment is not disliked or does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Comment successfully undisliked", http.StatusOK, "success")
}

func Temp_getDislikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	likes, err := database.Temp_selectDislikes()
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(likes)
}
