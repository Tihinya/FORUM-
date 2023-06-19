package controllers

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/router"
	"net/http"
	"time"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postId, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	if len(comment.Content) == 0 {
		http.Error(w, "Comment creation failed, the comment content can not be empty", http.StatusBadRequest)
		return
	}

	comment.CreationDate = time.Now()

	if !database.CreateCommentRow(comment, postId) {
		fmt.Fprintf(w, "Comment creation failed, post with id %v does not exist", postId)
		return
	}
	fmt.Fprint(w, "Comment successfully created")
}

func ReadComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	comment := database.SelectComment(commentId)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", comment)
}

func ReadComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postId, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	comment := database.SelectAllComments(postId)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", comment)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment

	commentId, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !database.UpdateComment(comment, commentId) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Comment updating failed, the comment with id %v does not exist", commentId)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Comment successfully updated")
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	if !database.DeleteComment(commentId) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Comment deletion failed, the comment with id %v does not exist", commentId)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Comment successfully deleted")
}
