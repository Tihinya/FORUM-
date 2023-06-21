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
		ReturnMessageJSON(w, r, "Invalid request body", 400, "error")
		return
	}

	postId, err := router.GetFieldString(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(comment.Content) == 0 {
		ReturnMessageJSON(w, r, "Comment creation failed, the comment content can not be empty", 400, "error")
		return
	}

	comment.CreationDate = time.Now()

	if !database.CreateCommentRow(comment, postId) {
		ReturnMessageJSON(w, r, "Comment creation failed, post with given ID does not exist", 400, "error")
		return
	}

	ReturnMessageJSON(w, r, "Comment successfully created", 200, "success")
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

	w.WriteHeader(http.StatusOK)
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment

	commentId, err := router.GetFieldString(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ReturnMessageJSON(w, r, "Invalid request body", 400, "error")
		return
	}

	if !database.UpdateComment(comment, commentId) {
		ReturnMessageJSON(w, r, "Comment updating failed, the comment with given id does not exist", 400, "error")
		return
	}

	w.WriteHeader(http.StatusOK)
	ReturnMessageJSON(w, r, "Comment successfully updated", 200, "success")
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := router.GetFieldString(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !database.DeleteComment(commentId) {
		ReturnMessageJSON(w, r, "Comment updating failed, the comment with given id does not exist", 400, "error")
		return
	}

	w.WriteHeader(http.StatusOK)
	ReturnMessageJSON(w, r, "Comment successfully deleted", 200, "success")
}
