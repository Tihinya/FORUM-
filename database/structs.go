package database

import "time"

type Post struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	UserInfo     UserInfo  `json:"userInfo"`
	CreationDate time.Time `json:"creationDate"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	Comments     string    `json:"comments"`
	Categories   []string  `json:"categories"`
	LastEdited   time.Time `json:"lastEdited"`
}

type Comment struct {
	Id           int       `json:"id"`
	Content      string    `json:"content"`
	UserInfo     UserInfo  `json:"userInfo"`
	CreationDate time.Time `json:"creationDate"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	LastEdited   time.Time `json:"lastEdited"`
	PostId       int       `json:"postId"`
}

type UserInfo struct {
	ID                   int    `json:"id"`
	ProfilePicture       string `json:"avatar"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type Like struct {
	Id        int    `json:"id"`
	PostId    int    `json:"postId"`
	CommentId int    `json:"commentId"`
	Username  string `json:"username"`
}
