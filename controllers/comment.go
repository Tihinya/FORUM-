package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/router"
	"log"
	"net/http"
	"time"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		returnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(comment.Content) == 0 {
		returnMessageJSON(w, "Comment creation failed, the comment content can not be empty", http.StatusBadRequest, "error")
		return
	}

	comment.CreationDate = time.Now()

	if !database.CreateCommentRow(comment, postId) {
		returnMessageJSON(w, "Comment creation failed, post with given ID does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully created", http.StatusOK, "success")
}

func ReadComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := router.GetFieldString(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment := database.SelectComment(commentId)

	json.NewEncoder(w).Encode(comment)
}

func ReadComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postId, err := router.GetFieldString(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comments := database.SelectAllComments(postId)

	json.NewEncoder(w).Encode(comments)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		returnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	if !database.UpdateComment(comment, commentId) {
		returnMessageJSON(w, "Comment updating failed, the comment with given id does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully updated", http.StatusOK, "success")
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !database.DeleteComment(commentId) {
		returnMessageJSON(w, "Comment updating failed, the comment with given id does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully deleted", http.StatusOK, "success")
}
