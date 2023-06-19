package controllers

import (
	"encoding/json"
	"fmt"
	"forum/database"
	"forum/router"
	"net/http"
	"time"
)

// Posts are readable on https://localhost:8080/posts

// POST method
func CreatePost(w http.ResponseWriter, r *http.Request) {

	var post database.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(post.Title) == 0 || len(post.Content) == 0 {
		http.Error(w, "Post creation failed, the post content or title can not be empty", http.StatusBadRequest)
		return
	}

	post.CreationDate = time.Now()

	database.CreatePost(post)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post successfully created")
}

// GET method
func ReadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postID, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	post := database.SelectPost(postID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", post)
}

// GET method for all posts
func ReadPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	posts := database.SelectAllPosts()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", posts)
}

// PATCH method
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post database.Post
	var exists bool

	postID, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(post.Content) == 0 {
		http.Error(w, "Post updating failed, the post content cannot be empty", http.StatusBadRequest)
		return
	}

	exists = database.UpdatePost(post, postID)

	if !exists {
		http.Error(w, "Post updating failed, the post with that ID does not exist", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post successfully updated")
}

// DELETE method
func DeletePost(w http.ResponseWriter, r *http.Request) {
	var exists bool

	postID, err := router.GetFieldString(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	exists = database.DeletePost(postID)

	if !exists {
		http.Error(w, "Post deletion failed, the post with that ID does not exist", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post successfully deleted")
}
