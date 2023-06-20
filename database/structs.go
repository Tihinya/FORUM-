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
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
}
