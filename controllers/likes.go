package controllers

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/router"
	"log"
	"net/http"
)

var username = "brozkie2"

func LikePost(w http.ResponseWriter, r *http.Request) {
	var existsLiked bool

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	// Authentication here maybe
	existsLiked, err = database.LikePost(postId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Liking post failed, post with id %v is already liked or does not exist", postId)
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

	// Authentication here maybe

	existsLiked, err = database.UnlikePost(postId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unliking post failed, post with id %v is not liked or does not exist", postId)
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

	// Authentication here maybe
	existsLiked, err = database.LikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Liking comment failed, comment with id %v is already liked or does not exist", commentId)
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

	// Authentication here maybe
	existsLiked, err = database.UnlikeComment(commentId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	if !existsLiked {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unliking comment failed, comment with id %v is not liked or does not exist", commentId)
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
