package database

import "time"

type Post struct {
	Id           int          `json:"id"`
	Title        string       `json:"title"`
	Content      string       `json:"content"`
	Image        string       `json:"image"`
	UserInfo     PostUserInfo `json:"user_info"`
	CreationDate time.Time    `json:"creation_date"`
	Likes        int          `json:"likes"`
	Dislikes     int          `json:"dislikes"`
	CommentCount int          `json:"comment_count"`
	Comments     string       `json:"comments"`
	Categories   []string     `json:"categories"`
	LastEdited   *time.Time   `json:"last_edited"`
}

type Category struct {
	Id       int    `json:"id"`
	Category string `json:"category"`
}

type PostCategory struct {
	PostId     int
	CategoryId int
}

type Comment struct {
	Id           int          `json:"id"`
	Content      string       `json:"content"`
	UserInfo     PostUserInfo `json:"user_info"`
	CreationDate time.Time    `json:"creation_date"`
	Likes        int          `json:"likes"`
	Dislikes     int          `json:"dislikes"`
	LastEdited   *time.Time   `json:"last_edited"`
	PostId       int          `json:"post_id"`
}

type PostUserInfo struct {
	ProfilePicture string `json:"avatar"`
	Username       string `json:"username"`
}

type UserInfo struct {
	ID                   int    `json:"id,omitempty"`
	ProfilePicture       string `json:"avatar,omitempty"`
	Username             string `json:"username,omitempty"`
	Email                string `json:"email,omitempty"`
	Password             string `json:"password,omitempty"`
	PasswordConfirmation string `json:"password_confirmation,omitempty"`
	RoleID               int    `json:"role_id,omitempty"`
}
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	ID      int    `json:"id"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"roleName,omitempty"`
}

type Like struct {
	PostId    int    `json:"post_id"`
	CommentId int    `json:"comment_id"`
	Username  string `json:"username"`
}

type Dislike struct {
	PostId    int    `json:"post_id"`
	CommentId int    `json:"comment_id"`
	Username  string `json:"username"`
}
