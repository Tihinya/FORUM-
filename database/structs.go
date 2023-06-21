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
	ID       int    `json:"id"`
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegistrationRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
	Avatar               string `json:"avatar"`
}
