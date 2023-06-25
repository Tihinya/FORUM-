package database

import (
	"fmt"
	"time"
)

func CreatePost(post Post) {
	stmt, err := db.Prepare(`
		INSERT INTO post (
			title,
			content,
			profile_picture,
			username,
			creation_date
		) VALUES (?, ?, ?, ?, ?)
	`)
	checkErr(err)

	result, err := stmt.Exec(post.Title, post.Content, post.UserInfo.ProfilePicture, post.UserInfo.Username, post.CreationDate)
	checkErr(err)

	postId, err := result.LastInsertId()
	checkErr(err)

	addCategory(post, int(postId))
}

func SelectPost(id string) []Post {
	var posts []Post

	rows, err := db.Query(`
		SELECT post.id, post.title, post.content,
		post.profile_picture, post.username, post.creation_date,
		post.likes, post.dislikes, post.last_edited
		FROM post WHERE id = ?
	`, id)
	checkErr(err)

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
		)
		checkErr(err)

		post.Categories = getCategories(post)
		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)
		posts = append(posts, post)
	}

	return posts
}

// GET all posts from posts table
func SelectAllPosts() []Post {
	var posts []Post
	rows, err := db.Query(`
		SELECT post.id, post.title, post.content,
		post.profile_picture, post.username, post.creation_date,
		post.likes, post.dislikes, post.last_edited
		FROM post
	`)

	checkErr(err)

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
		)

		post.Categories = getCategories(post)

		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)

		posts = append(posts, post)
	}

	return posts
}

func DeletePost(postId int) bool {
	var post Post

	if !checkIfPostExist(postId) {
		return false
	}

	post.Categories = nil
	UpdatePost(post, postId)

	stmt, err := db.Prepare(`
		DELETE FROM post WHERE id = ?
	`)
	checkErr(err)

	_, err = stmt.Exec(postId)
	checkErr(err)

	return true
}

func UpdatePost(post Post, postID int) bool {
	if !checkIfPostExist(postID) {
		return false
	}

	stmt, err := db.Prepare(`
		UPDATE post SET
			title = ?,
			content = ?,
			last_edited = ?
		WHERE id = ?
	`)
	checkErr(err)

	if post.Id != 0 {
		postID = int(post.Id)
	}

	_, err = stmt.Exec(post.Title, post.Content, time.Now(), postID)
	checkErr(err)

	updateCategories(post, postID)

	return true
}

func getCategories(post Post) []string {

	categoryRows, err := db.Query(`
		SELECT category FROM category
		INNER JOIN post_category ON category.id = post_category.category_id
		INNER JOIN post ON post_category.post_id = post.id
		WHERE post.id = ?
	`, post.Id)
	checkErr(err)

	for categoryRows.Next() {
		var category string

		err = categoryRows.Scan(&category)
		checkErr(err)

		post.Categories = append(post.Categories, category)
	}

	return post.Categories
}

func updateCategories(post Post, postId int) {
	var existingCategories []int64

	rows, err := db.Query(`
		SELECT category_id FROM post_category
		WHERE post_id = ?
	`, postId)
	checkErr(err)

	for rows.Next() {
		var categoryId int64

		err = rows.Scan(&categoryId)
		checkErr(err)

		existingCategories = append(existingCategories, categoryId)
	}

	// If no categories in current post, skip straight to adding categories
	if existingCategories == nil {
		addCategory(post, postId)
		return
	}

	for _, categoryId := range existingCategories {
		var count int

		// Get the count of categories linked to a post
		err := db.QueryRow(`
			SELECT COUNT(*) FROM post_category
			WHERE category_id = ?
		`, categoryId).Scan(&count)
		checkErr(err)

		stmt, err := db.Prepare(`
			DELETE FROM post_category WHERE category_id = ? AND post_id = ?
		`)
		checkErr(err)

		_, err = stmt.Exec(categoryId, postId)
		checkErr(err)

		// Removes all categories with only 1 connection to posts
		if count == 1 {
			stmt, err := db.Prepare(`
				DELETE FROM category WHERE id = ?
			`)
			checkErr(err)

			_, err = stmt.Exec(categoryId)
			checkErr(err)

		}
	}

	addCategory(post, postId)
}

// Categories many-to-many link
func addCategory(post Post, postId int) {
	var postCount int

	for i := range post.Categories {
		var categoryId int64

		err := db.QueryRow(`SELECT COUNT(*) FROM post`).Scan(&postCount)
		checkErr(err)

		if checkIfCategoryExist(post.Categories[i]) {
			err = db.QueryRow("SELECT id FROM category WHERE category = ?", post.Categories[i]).Scan(&categoryId)

		} else {
			stmt, err := db.Prepare(`INSERT INTO category (category) VALUES (?)`)
			checkErr(err)

			result, err := stmt.Exec(post.Categories[i])
			checkErr(err)

			categoryId, err = result.LastInsertId()
			checkErr(err)

			if postCount == 1 {
				insertPostCategory(postId, categoryId)
			}
		}

		if postCount > 1 {
			insertPostCategory(postId, categoryId)
		}

	}
}

func insertPostCategory(postId int, categoryId int64) {
	stmt, err := db.Prepare(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`)
	checkErr(err)

	_, err = stmt.Exec(postId, categoryId)
	checkErr(err)
}

func SelectAllCategories() []Category {
	var categories []Category
	rows, err := db.Query("SELECT * FROM category")
	checkErr(err)

	for rows.Next() {
		var category Category

		err = rows.Scan(&category.Id, &category.Category)
		checkErr(err)

		categories = append(categories, category)
	}

	return categories
}

func SelectAllPostCategory() []PostCategory {
	var post_categories []PostCategory
	rows, err := db.Query("SELECT * FROM post_category")
	checkErr(err)

	for rows.Next() {
		var post_category PostCategory

		err = rows.Scan(&post_category.PostId, &post_category.CategoryId)
		checkErr(err)

		post_categories = append(post_categories, post_category)
	}

	return post_categories
}

func checkIfPostExist(commentId int) bool {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE id=?)", commentId).Scan(&exists)

	return err == nil && exists
}

func checkIfCategoryExist(category string) bool {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM category WHERE category=?)", category).Scan(&exists)

	return err == nil && exists
}

func getCategoryId(category string) int64 {
	var categoryId int64

	err := db.QueryRow("SELECT 1 FROM category WHERE id = ?", category).Scan(&categoryId)

	if err == nil {
		return categoryId
	}

	return 0
}
