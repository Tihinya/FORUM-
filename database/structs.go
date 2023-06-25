package database

import "time"

type Post struct {
	Id           int        `json:"id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	UserInfo     UserInfo   `json:"user_info"`
	CreationDate time.Time  `json:"creation_date"`
	Likes        int        `json:"likes"`
	Dislikes     int        `json:"dislikes"`
	Comments     string     `json:"comments"`
	Categories   []string   `json:"categories"`
	LastEdited   *time.Time `json:"last_edited"`
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
	Id           int        `json:"id"`
	Content      string     `json:"content"`
	UserInfo     UserInfo   `json:"user_info"`
	CreationDate time.Time  `json:"creation_date"`
	Likes        int        `json:"likes"`
	Dislikes     int        `json:"dislikes"`
	LastEdited   *time.Time `json:"last_edited"`
	PostId       int        `json:"post_id"`
}

type UserInfo struct {
	ProfilePicture string `json:"profile_picture"`
	Username       string `json:"username"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
