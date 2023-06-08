package controllers

import (
	"encoding/json"
	"forum/router"
	"net/http"
	"time"
)

type Post struct {
	PostId       int       `json:"postId"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	UserInfo     UserInfo  `json:"userInfo"`
	CreationDate time.Time `json:"creationDate"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	Categories   []string  `json:"categories"`
}

type UserInfo struct {
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
}

var tempDB = make(map[int]Post)

func CreatePost(w http.ResponseWriter, r *http.Request) {

	var post Post

	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err2 := json.NewDecoder(r.Body).Decode(&post)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, exist := tempDB[postID]

	if !exist {
		post.CreationDate = time.Now()
		post.Likes = 0
		post.Dislikes = 0
		post.PostId = postID

		// Send data to database
		tempDB[postID] = post

		json.NewEncoder(w).Encode("Post successfully created")
	} else {
		json.NewEncoder(w).Encode("Post creation failed, a post on that ID already exists")
	}

}
func ReadPost(w http.ResponseWriter, r *http.Request) {
	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}
	post, exist := tempDB[postID]

	if !exist {
		http.Error(w, "Post does not exist, failed to GET", http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post Post

	postID, err := router.GetFieldInt(r, "id")
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}

	err2 := json.NewDecoder(r.Body).Decode(&post)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldPost, exist := tempDB[postID]

	if !exist {
		http.Error(w, "Post does not exist, failed to update", http.StatusBadRequest)
	} else {
		post.CreationDate = oldPost.CreationDate
		post.Likes = oldPost.Likes
		post.Dislikes = oldPost.Dislikes
		post.PostId = postID
		tempDB[postID] = post

		json.NewEncoder(w).Encode("Post successfully updated")
	}

}
func DeletePost(w http.ResponseWriter, r *http.Request) {
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
}

/* Test curl:

curl -X POST -H "Content-Type: application/json" -d '{
  "title": "Test post",
  "text": "This is the content of my post",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "john_doe"
  },
  "likes": 0,
  "dislikes": 0,
  "categories": ["technology", "programming"]
}' -k https://localhost:8080/post/1

curl -X POST -H "Content-Type: application/json" -d '{
  "title": "2222222222",
  "text": "2 This is the content of my post 2",
  "userInfo": {
    "avatar": "https://example.com/avatar.png",
    "username": "john_doe"
  },
  "categories": ["technology", "programming"]
}' -k https://localhost:8080/post/2

curl -X GET -k https://localhost:8080/post/1

curl -X PATCH -H "Content-Type: application/json" -d '{
  "title": "UPDATED UPDATED UPDATED",
  "content": "Updated Updated Updated?",
  "categories": ["updated", "the whats?"]
}' -k https://localhost:8080/post/1

curl -X DELETE -k https://localhost:8080/post/1

*/
