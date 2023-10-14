package database

import (
	"fmt"
	"strings"
	"time"
)

func CreatePost(post Post) error {
	stmt, err := DB.Prepare(`
		INSERT INTO post (
			title,
			content,
			profile_picture,
			username,
			image,
			creation_date
		) VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(post.Title, post.Content, post.UserInfo.ProfilePicture, post.UserInfo.Username, post.Image, post.CreationDate)
	if err != nil {
		return err
	}

	postId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	err = addCategory(post, int(postId))
	if err != nil {
		return err
	}

	return nil
}

func SelectPost(id string) ([]Post, error) {
	posts := make([]Post, 0)

	rows, err := DB.Query(`
		SELECT post.id, post.title, post.content,
		post.profile_picture, post.username, post.creation_date,
		post.likes, post.dislikes, post.last_edited, post.image
		FROM post WHERE id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.UserInfo.ProfilePicture,
			&post.UserInfo.Username,
			&post.CreationDate,
			&post.Likes,
			&post.Dislikes,
			&post.LastEdited,
			&post.Image,
		)
		if err != nil {
			return nil, err
		}

		post.Categories, err = getCategories(post)
		if err != nil {
			return nil, err
		}

		post.Likes, _ = getPostLikes(post.Id)
		post.Dislikes, _ = getPostDislikes(post.Id)

		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)
		post.CommentCount = getCommentsCount(post.Id)
		post.UserInfo.ProfilePicture, _ = GetAvatar(post.UserInfo.Username)

		posts = append(posts, post)
	}

	return posts, nil
}

// GET all posts from posts table
func SelectAllPosts(categoriesString string) ([]Post, error) {
	posts := make([]Post, 0)

	rows, err := DB.Query(`
		SELECT post.id, post.title, post.content,
		post.profile_picture, post.username, post.creation_date,
		post.likes, post.dislikes, post.last_edited, post.image
		FROM post
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.UserInfo.ProfilePicture,
			&post.UserInfo.Username,
			&post.CreationDate,
			&post.Likes,
			&post.Dislikes,
			&post.LastEdited,
			&post.Image,
		)
		if err != nil {
			return nil, err
		}

		post.Categories, err = getCategories(post)
		if err != nil {
			return nil, err
		}

		post.Likes, _ = getPostLikes(post.Id)
		post.Dislikes, _ = getPostDislikes(post.Id)

		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)
		post.CommentCount = getCommentsCount(post.Id)
		post.UserInfo.ProfilePicture, _ = GetAvatar(post.UserInfo.Username)

		if contains(post.Categories, categoriesString) {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func DeletePost(postId int, username string) (bool, error) {
	var post Post
	var exists bool

	if !checkPostOwnership(postId, username) {
		return false, nil
	}

	if !checkIfPostExist(postId) {
		return false, nil
	}

	// For deleting leftover categories
	post.Categories = nil
	exists, err = UpdatePost(post, postId, username)
	if !exists {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare(`
		DELETE FROM post WHERE id = ?
	`)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(postId)
	if err != nil {
		return false, err
	}

	commentIds, err := getPostCommentIds(postId)
	if err != nil {
		return false, err
	}

	err = deletePostComments(postId)
	if err != nil {
		return false, err
	}

	err = deletePostLikes(postId)
	if err != nil {
		return false, err
	}

	for _, commentId := range commentIds {
		err = deleteCommentLikes(commentId)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func DeletePostModerator(postId int) (bool, error) {
	var post Post
	var exists bool

	if !checkIfPostExist(postId) {
		return false, nil
	}

	// For deleting leftover categories
	post.Categories = nil
	exists, err = UpdatePost(post, postId, "admin")
	if !exists {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare(`
		DELETE FROM post WHERE id = ?
	`)
	if err != nil {
		return false, err
	}

	_, err = stmt.Exec(postId)
	if err != nil {
		return false, err
	}

	commentIds, err := getPostCommentIds(postId)
	if err != nil {
		return false, err
	}

	err = deletePostComments(postId)
	if err != nil {
		return false, err
	}

	err = deletePostLikes(postId)
	if err != nil {
		return false, err
	}

	for _, commentId := range commentIds {
		err = deleteCommentLikes(commentId)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func UpdatePost(post Post, postID int, username string) (bool, error) {
	if !checkPostOwnership(postID, username) && username != "admin" {
		return false, nil
	}

	if !checkIfPostExist(postID) {
		return false, nil
	}

	stmt, err := DB.Prepare(`
		UPDATE post SET
			title = ?,
			content = ?,
			last_edited = ?
		WHERE id = ?
	`)
	if err != nil {
		return false, err
	}

	if post.Id != 0 {
		postID = int(post.Id)
	}

	_, err = stmt.Exec(post.Title, post.Content, time.Now(), postID)
	if err != nil {
		return false, err
	}

	err = updateCategories(post, postID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func deletePostLikes(postId int) error {
	stmt, err := DB.Prepare(`
		DELETE FROM like WHERE post_id = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(postId)
	if err != nil {
		return err
	}

	stmt, err = DB.Prepare(`
		DELETE FROM dislike WHERE post_id = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(postId)
	if err != nil {
		return err
	}

	return nil
}

func getPostCreatorByPostId(postId int) (string, error) {
	var username string

	err := DB.QueryRow("SELECT username FROM post WHERE id = ?", postId).Scan(&username)

	if err == nil {
		return username, nil
	}

	return "error", err
}

func checkIfPostExist(commentId int) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE id = ?)", commentId).Scan(&exists)

	return err == nil && exists
}

func checkPostOwnership(postId int, username string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE id = ? AND username = ?)", postId, username).Scan(&exists)

	return err == nil && exists
}

func contains(postArr []string, urlParams string) bool {
	var found bool

	urlCategories := strings.Split(urlParams, ",")

	if len(urlParams) == 0 {
		return true
	}

	for j := range urlCategories {
		for i := range postArr {
			if urlCategories[j] == postArr[i] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
		found = false
	}

	return true
}
