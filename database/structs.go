package database

import "time"

type Post struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"text"`
	UserInfo     UserInfo  `json:"userInfo"`
	CreationDate time.Time `json:"creationDate"`
	Likes        int       `json:"likes"`
	Dislikes     int       `json:"dislikes"`
	Categories   []string  `json:"categories"`
	LastEdited   time.Time `json:"lastEdited"`
}

type UserInfo struct {
	ID                   int    `json:"id"`
	ProfilePicture       string `json:"avatar"`
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password,omitempty"`
	PasswordConfirmation string `json:"password_confirmation,omitempty"`
}
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
