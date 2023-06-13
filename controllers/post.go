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

	if len(post.Title) == 0 || len(post.Text) == 0 {
		http.Error(w, "Post creation failed, the post content can not be empty", http.StatusBadRequest)
		return
	}

	post.CreationDate = time.Now()

	database.CreatePost(post)
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
	/*var post database.Post

	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	err2 := json.NewDecoder(r.Body).Decode(&post)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(post.Text) == 0 || len(post.Title) == 0 {
		http.Error(w, "Post updating failed, the post content can not be empty", http.StatusBadRequest)
		return
	}

	oldPost, exist := tempDB[postID]

	if !exist {
		http.Error(w, "Post does not exist, failed to update", http.StatusBadRequest)
	} else {
		post.CreationDate = oldPost.CreationDate
		post.Likes = oldPost.Likes
		post.Dislikes = oldPost.Dislikes
		post.Id = postID
		post.UserInfo.Avatar = oldPost.UserInfo.Avatar
		post.UserInfo.Username = oldPost.UserInfo.Username
		tempDB[postID] = post

		json.NewEncoder(w).Encode("Post successfully updated")
	}
	*/
}

// DELETE Method
func DeletePost(w http.ResponseWriter, r *http.Request) {
	/*
		postID, err := router.GetFieldInt(r, "id")
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
		}

		_, exist := tempDB[postID]

		if !exist {
			http.Error(w, "Post does not exist, failed to delete", http.StatusBadRequest)
		} else {
			delete(tempDB, postID)

			json.NewEncoder(w).Encode("Post successfully deleted")
		}
	*/
}
