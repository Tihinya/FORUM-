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
		http.Error(w, "Post creation failed, the post content can not be empty", http.StatusBadRequest)
		return
	}

	post.CreationDate = time.Now()

	database.CreatePost(post)
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

	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.UpdatePost(post, postID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post successfully updated")
}

// DELETE method
func DeletePost(w http.ResponseWriter, r *http.Request) {

	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	database.DeletePost(postID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Post successfully deleted")
}
