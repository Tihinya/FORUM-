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

	if len(comment.Content) == 0 {
		http.Error(w, "Comment creation failed, the comment content can not be empty", http.StatusBadRequest)
		return
	}

	comment.CreationDate = time.Now()

	database.CreateCommentRow(comment)
	fmt.Fprint(w, "Comment successfully created")
}
func ReadComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentID, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	comment := database.SelectComment(commentID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", comment)
}
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment

	commentID, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.UpdateComment(comment, commentID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Comment successfully updated")
}
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
	}

	database.DeleteComment(commentID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Comment successfully deleted")
}
