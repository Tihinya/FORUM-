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
	var exists bool

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

	exists, err = database.CreateCommentRow(comment, postId)

	comment.CreationDate = time.Now()

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal error", http.StatusBadRequest, "error")
		return
	}
	if !exists {
		returnMessageJSON(w, "Comment creation failed, post with given ID does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully created", http.StatusOK, "success")
}

func ReadComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment, err := database.SelectComment(commentId)
	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(comment)
}

func ReadComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comments, err := database.SelectAllComments(postId)
	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment
	var exists bool

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

	exists, err = database.UpdateComment(comment, commentId)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		returnMessageJSON(w, "Comment updating failed, the comment with given id does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully updated", http.StatusOK, "success")
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var exists bool

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	exists, err = database.DeleteComment(commentId)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		returnMessageJSON(w, "Comment updating failed, the comment with given id does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully deleted", http.StatusOK, "success")
}
