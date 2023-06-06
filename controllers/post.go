package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {

	// Get post data via formvalue(current) or pure json?
	postTitle := r.FormValue("post__title")
	postContent := r.FormValue("post__content")
	postCategories := r.FormValue("post__tags")
	postCreator := r.FormValue("username")
	postCreationDate := time.Now().String()

	// Parse data to json
	data := map[string]string{
		"title":        postTitle,
		"content":      postContent,
		"categories":   postCategories,
		"creator":      postCreator,
		"creationDate": postCreationDate,
	}
	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))

	// Send data to database
	// db.CreatePost(jsonData)

}
func ReadPost(w http.ResponseWriter, r *http.Request) {} // Is this even necessary?

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	postTitle := r.FormValue("post__title")
	postContent := r.FormValue("post__content")
	postCategories := r.FormValue("post__tags")
	postCreator := r.FormValue("username")
	postCreationDate := time.Now().String()
	postId := r.FormValue("post__id")

	// Parse data to json
	data := map[string]string{
		"title":        postTitle,
		"content":      postContent,
		"categories":   postCategories,
		"creator":      postCreator,
		"creationDate": postCreationDate,
		"postId":       postId,
	}
	jsonData, _ := json.Marshal(data)
	fmt.Println(string(jsonData))

	// Send data to database
	// db.EditPost(jsonData)
}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	postId := r.FormValue("post__id")
	fmt.Println(postId)
	// Delete from database
	// db.DeletePost(postId)
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

*/
