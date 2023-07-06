package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/router"
	"log"
	"net/http"
	"time"
)

// Posts are readable on https://localhost:8080/posts

// POST method
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post database.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		returnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	if len(post.Title) == 0 || len(post.Content) == 0 {
		returnMessageJSON(w, "Post creation failed, the post content or title can not be empty", http.StatusBadRequest, "error")
		return
	}

	post.CreationDate = time.Now()

	err = database.CreatePost(post)
	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	returnMessageJSON(w, "Post successfully created", http.StatusOK, "success")
}

// GET method
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
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(post)
}

// GET method for all posts
func ReadPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	posts, err := database.SelectAllPosts()
	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(posts)
}

// PATCH method
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
		returnMessageJSON(w, "Invalid request body", http.StatusBadRequest, "error")
		return
	}

	if len(post.Content) == 0 {
		returnMessageJSON(w, "Post updating failed, the post content cannot be empty", http.StatusBadRequest, "error")
		return
	}

	exists, err = database.UpdatePost(post, postID)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		returnMessageJSON(w, "Post updating failed, the post with that ID does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Post successfully updated", http.StatusOK, "success")
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

	exists, err = database.DeletePost(postID)

	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}
	if !exists {
		returnMessageJSON(w, "Post deletion failed, the post with that ID does not exist", http.StatusBadRequest, "error")
		return
	}

	returnMessageJSON(w, "Post successfully deleted", http.StatusOK, "success")
}

func ReadCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := database.SelectAllCategories()
	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(categories)
}

func ReadPostCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	post_categories, err := database.SelectAllPostCategory()
	if err != nil {
		log.Println(err)
		returnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(post_categories)
}
