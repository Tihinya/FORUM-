package controllers

import (
	"encoding/json"
	"fmt"
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
		ReturnErrorMessageJSON(w, r, "Invalid request body", 400)
		return
	}

	postId, err := router.GetFieldString(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(comment.Content) == 0 {
		ReturnErrorMessageJSON(w, r, "Comment creation failed, the comment content can not be empty", 400)
		return
	}

	comment.CreationDate = time.Now()

	if !database.CreateCommentRow(comment, postId) {
		ReturnErrorMessageJSON(w, r, "Comment creation failed, post with given ID does not exist", 400)
		return
	}

	fmt.Fprint(w, "Comment successfully created")
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
		ReturnErrorMessageJSON(w, r, "Invalid request body", 400)
		return
	}

	if !database.UpdateComment(comment, commentId) {
		ReturnErrorMessageJSON(w, r, "Comment updating failed, the comment with given id does not exist", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Comment successfully updated")
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := router.GetFieldString(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !database.DeleteComment(commentId) {
		ReturnErrorMessageJSON(w, r, "Comment updating failed, the comment with given id does not exist", 400)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Comment successfully deleted")
}
