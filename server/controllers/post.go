package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"forum/database"
	"forum/router"
	"forum/session"
	"forum/validation"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post database.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	if len(post.Title) == 0 || len(post.Content) == 0 {
		ReturnMessageJSON(w, "Post creation failed, the post content or title can not be empty", http.StatusBadRequest, "error")
		return
	}

	if !containsStr(post.Categories, ',') {
		ReturnMessageJSON(w, "Post creation failed, the post categories cannot contain a comma", http.StatusBadRequest, "error")
		return
	}

	post.CreationDate = time.Now()
	post.UserInfo.Username = getUsername(r)
	post.UserInfo.ProfilePicture = "https://example.com/avatar.png"

	err = database.CreatePost(post)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "Post successfully created", http.StatusOK, "success")
}

func ReadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postID, err := router.GetFieldString(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	post, err := database.SelectPost(postID)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(post)
}

func ReadPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories := r.URL.Query().Get("categories")

	posts, err := database.SelectAllPosts(categories)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post database.Post
	var exists bool

	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	if len(post.Content) == 0 {
		ReturnMessageJSON(w, "Post updating failed, the post content cannot be empty", http.StatusBadRequest, "error")
		return
	}

	if !containsStr(post.Categories, ',') {
		ReturnMessageJSON(w, "Post updating failed, the post categories cannot contain a comma", http.StatusBadRequest, "error")
		return
	}

	username := getUsername(r)
	exists, err = database.UpdatePost(post, postID, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		ReturnMessageJSON(w, "Post updating failed, you do not have the right to update this post or the post with that ID does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Post successfully updated", http.StatusOK, "success")
}

// DELETE method
func DeletePost(w http.ResponseWriter, r *http.Request) {
	var exists bool

	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := getUsername(r)
	exists, err = database.DeletePost(postID, username)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		ReturnMessageJSON(w, "Post deletion failed, you do not have the right to delete this post or the post with that ID does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Post successfully deleted", http.StatusOK, "success")
}

func DeletePostModerator(w http.ResponseWriter, r *http.Request) {
	var exists bool

	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	exists, err = database.DeletePostModerator(postID)

	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		ReturnMessageJSON(w, "Post deletion failed, the post with that ID does not exist", http.StatusBadRequest, "error")
		return
	}

	ReturnMessageJSON(w, "Post successfully deleted", http.StatusOK, "success")
}

func CreateReportPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var postReport database.PostReport
	err := json.NewDecoder(r.Body).Decode(&postReport)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	// Checking if there's a post
	exist, err := validation.HasPost(database.DB, postReport.PostID)
	if err != nil {
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exist {
		ReturnMessageJSON(w, "There's no such post", http.StatusBadRequest, "error")
		return
	}
	// Сhecking  if there's a  Post Report
	exist, err = validation.HasPendingPostReport(database.DB, postReport.PostID)
	if err != nil {
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if exist {
		ReturnMessageJSON(w, "This post has already been reported and is awaiting a decision ", http.StatusBadRequest, "error")
		return
	}

	//Get User ID
	session, err := session.SessionStorage.GetSession(r)
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}
	err = database.CreatePostReport(postReport, session.UserId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "Report post successfully created", http.StatusOK, "success")
}

func ReadReportPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	postReport, err := database.SelectAllPostReports()
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(postReport)

}

func ReadReportPostAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	session, err := session.SessionStorage.GetSession(r)
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}
	// Check if a valid session exists
	if session == nil {
		ReturnMessageJSON(w, "Unauthorized", http.StatusUnauthorized, "error")
		return
	}

	postReports, err := database.GetUserReportedPost(session.UserId)
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(postReports)

}

func UpdatedReportPost(w http.ResponseWriter, r *http.Request) {
	var postReport database.PostReport
	err := json.NewDecoder(r.Body).Decode(&postReport)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	if postReport.Status == "" {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	if postReport.Status != "approved" && postReport.Status != "rejected" {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	// Give Report Post Status
	status, err := validation.GetReportPostStatus(database.DB, postReport.ReportID)
	if err != nil {
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if status != "pending" {
		ReturnMessageJSON(w, "The resolution has already been approved", http.StatusBadRequest, "error")
		return
	}

	// Update User Role
	err = database.UpdatePostReport(postReport)
	if err != nil {
		ReturnMessageJSON(w, "Report Post was updated", http.StatusInternalServerError, "error")
		return
	}
	// Delete RoleRequest
	ReturnMessageJSON(w, "Application approved", http.StatusOK, "success")
}
func UpdatedReportPostAnswer(w http.ResponseWriter, r *http.Request) {
	var postReport database.PostReport
	err := json.NewDecoder(r.Body).Decode(&postReport)
	if err != nil {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	if !postReport.Seen {
		ReturnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}
	// Update User Role
	err = database.UpdatePostReportAnswer(postReport)
	if err != nil {
		ReturnMessageJSON(w, "Internal Server Error", http.StatusInternalServerError, "error")
		return
	}
	// Delete RoleRequest
	ReturnMessageJSON(w, "Success, answer came back", http.StatusOK, "success")
}
