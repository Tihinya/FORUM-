package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"forum/database"
	"forum/router"
	"forum/session"
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

	// Authentication here
	sessionToken, sessionTokenFound := checkForSessionToken(r)
	if !sessionTokenFound {
		returnMessageJSON(w, "Session token not found", http.StatusUnauthorized, "unauthorized")
		return
	}

	if !checkIfUserLoggedin(sessionToken) {
		returnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "unauthorized")
		return
	}

	if len(comment.Content) == 0 {
		returnMessageJSON(w, "Comment creation failed, the comment content can not be empty", http.StatusBadRequest, "error")
		return
	}

	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	username, err := database.GetUsername(userID)
	if err != nil {
		returnMessageJSON(w, "You are not logged in", http.StatusInternalServerError, "unauthorized")
		return
	}

	comment.CreationDate = time.Now()
	comment.UserInfo.Username = username

	exists, err = database.CreateCommentRow(comment, postId)

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

	// Authentication here
	sessionToken, sessionTokenFound := checkForSessionToken(r)
	if !sessionTokenFound {
		returnMessageJSON(w, "Session token not found", http.StatusUnauthorized, "unauthorized")
		return
	}

	if !checkIfUserLoggedin(sessionToken) {
		returnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "unauthorized")
		return
	}

	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	username, err := database.GetUsername(userID)
	if err != nil {
		returnMessageJSON(w, "You are not logged in", http.StatusInternalServerError, "unauthorized")
		return
	}

	exists, err = database.UpdateComment(comment, commentId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		returnMessageJSON(w, "Comment updating failed, you do not have the right to update this comment or the comment with given id does not exist", http.StatusBadRequest, "error")
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

	// Authentication here
	sessionToken, sessionTokenFound := checkForSessionToken(r)
	if !sessionTokenFound {
		returnMessageJSON(w, "Session token not found", http.StatusUnauthorized, "unauthorized")
		return
	}

	if !checkIfUserLoggedin(sessionToken) {
		returnMessageJSON(w, "You are not logged in", http.StatusUnauthorized, "unauthorized")
		return
	}

	userID := session.SessionStorage.GetSession(sessionToken.Value).UserId
	username, err := database.GetUsername(userID)
	if err != nil {
		returnMessageJSON(w, "You are not logged in", http.StatusInternalServerError, "unauthorized")
		return
	}

	exists, err = database.DeleteComment(commentId, username)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		returnMessageJSON(w, "Comment updating failed, you do not have the right to update this comment or the comment with given id does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Comment successfully deleted", http.StatusOK, "success")
}
