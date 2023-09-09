package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/database"
	"forum/router"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	username := getUsername(r)

	// Check if post is already disliked
	if database.CheckIfPostDisliked(postId, username) {
		ReturnMessageJSON(w, "Failed to like post. Post is already disliked.", http.StatusBadRequest, "error")
		return
	}

	existsLiked, err = database.LikePost(postId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		ReturnMessageJSON(w, "Liking post failed, post is already liked or does not exist", http.StatusBadRequest, "error")
		return
	}

	err = database.CreateNotification("post", postId, "like")
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "Post successfully liked", http.StatusOK, "success")
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	username := getUsername(r)
	existsLiked, err = database.UnlikePost(postId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		ReturnMessageJSON(w, "Unliking post failed, post is not liked or does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Post successfully unliked", http.StatusOK, "success")
}

func LikeComment(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	username := getUsername(r)

	// Check if commennt is already disliked
	if database.CheckIfCommentDisliked(commentId, username) {
		ReturnMessageJSON(w, "Failed to like comment. Comment is already disliked.", http.StatusBadRequest, "error")
		return
	}

	existsLiked, err = database.LikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		ReturnMessageJSON(w, "Liking comment failed, comment is already liked or does not exist", http.StatusBadRequest, "error")
		return
	}

	err = database.CreateNotification("comment", commentId, "like")
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "Comment successfully liked", http.StatusOK, "success")
}

func UnlikeComment(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	username := getUsername(r)
	existsLiked, err = database.UnlikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		ReturnMessageJSON(w, "Unliking comment failed, comment is not liked or does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Comment successfully unliked", http.StatusOK, "success")
}

func Temp_getLikes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	likes, err := database.Temp_selectLikes()
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(likes)
}
