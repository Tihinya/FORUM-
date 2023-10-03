package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"forum/database"
	"forum/router"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment
	var exists bool

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	postId, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment.CreationDate = time.Now()
	comment.UserInfo.Username = getUsername(r)

	exists, err = database.CreateCommentRow(comment, postId, comment.UserInfo.Username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal error", http.StatusBadRequest, "error")
		return
	}
	if !exists {
		ReturnMessageJSON(w, "Comment creation failed, post with given ID does not exist", http.StatusBadRequest, "error")
		return
	}

	err = database.CreateNotification(postId, "post", postId, "comment")
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "Comment successfully created", http.StatusOK, "success")
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
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
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
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
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
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	username := getUsername(r)
	exists, err = database.UpdateComment(comment, commentId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		ReturnMessageJSON(w, "Comment updating failed, you do not have the right to update this comment or the comment with given id does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Comment successfully updated", http.StatusOK, "success")
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var exists bool

	commentId, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := getUsername(r)
	exists, err = database.DeleteComment(commentId, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		ReturnMessageJSON(w, "Comment updating failed, you do not have the right to update this comment or the comment with given id does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Comment successfully deleted", http.StatusOK, "success")
}
