package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"forum/database"
	"forum/router"
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

func ReadCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := database.SelectAllCategories()
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(categories)
}

func ReadPostCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postCategories, err := database.SelectAllPostCategory()
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(postCategories)
}
