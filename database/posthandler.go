package database

import (
	"fmt"
	"time"
)

func CreatePost(post Post) error {
	stmt, err := DB.Prepare(`
		INSERT INTO post (
			title,
			content,
			profile_picture,
			username,
			creation_date
		) VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(post.Title, post.Content, post.UserInfo.ProfilePicture, post.UserInfo.Username, post.CreationDate)
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
	var posts []Post

	rows, err := DB.Query(`
		SELECT post.id, post.title, post.content,
		post.profile_picture, post.username, post.creation_date,
		post.likes, post.dislikes, post.last_edited
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
		)
		if err != nil {
			return nil, err
		}

		post.Categories, err = getCategories(post)

		post.Likes, err = getPostLikes(post.Id)
		post.Dislikes, err = getPostDislikes(post.Id)

		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)

		posts = append(posts, post)
	}

	return posts, nil
}

// GET all posts from posts table
func SelectAllPosts() ([]Post, error) {
	var posts []Post

	rows, err := DB.Query(`
		SELECT post.id, post.title, post.content,
		post.profile_picture, post.username, post.creation_date,
		post.likes, post.dislikes, post.last_edited
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
		)
		if err != nil {
			return nil, err
		}

		post.Categories, err = getCategories(post)

		post.Likes, err = getPostLikes(post.Id)
		post.Dislikes, err = getPostDislikes(post.Id)

		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)

		posts = append(posts, post)
	}

	return posts, nil
}

func DeletePost(postId int) (bool, error) {
	var post Post
	var exists bool

	if !checkIfPostExist(postId) {
		return false, nil
	}

	post.Categories = nil
	exists, err = UpdatePost(post, postId)
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

	err = deletePostComments(postId)
	if err != nil {
		return false, err
	}

	err = deletePostLikes(postId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdatePost(post Post, postID int) (bool, error) {
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

func getCategories(post Post) ([]string, error) {
	categoryRows, err := DB.Query(`
		SELECT category FROM category
		INNER JOIN post_category ON category.id = post_category.category_id
		INNER JOIN post ON post_category.post_id = post.id
		WHERE post.id = ?
	`, post.Id)

	if err != nil {
		return nil, err
	}

	for categoryRows.Next() {
		var category string

		err = categoryRows.Scan(&category)
		if err != nil {
			return nil, err
		}

		post.Categories = append(post.Categories, category)
	}

	return post.Categories, nil
}

func updateCategories(post Post, postId int) error {
	var existingCategories []int64

	rows, err := DB.Query(`
		SELECT category_id FROM post_category
		WHERE post_id = ?
	`, postId)
	if err != nil {
		return err
	}

	for rows.Next() {
		var categoryId int64

		err = rows.Scan(&categoryId)
		if err != nil {
			return err
		}

		existingCategories = append(existingCategories, categoryId)
	}

	// If no categories in current post, skip straight to adding categories
	if existingCategories == nil {
		err = addCategory(post, postId)
		if err != nil {
			return err
		}
		return nil
	}

	for _, categoryId := range existingCategories {
		var count int

		// Get the count of categories linked to a post
		err := DB.QueryRow(`
			SELECT COUNT(*) FROM post_category
			WHERE category_id = ?
		`, categoryId).Scan(&count)
		if err != nil {
			return err
		}

		stmt, err := DB.Prepare(`
			DELETE FROM post_category WHERE category_id = ? AND post_id = ?
		`)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(categoryId, postId)
		if err != nil {
			return err
		}

		// Removes all categories with only 1 connection to posts
		if count == 1 {
			stmt, err := DB.Prepare(`
				DELETE FROM category WHERE id = ?
			`)
			if err != nil {
				return err
			}

			_, err = stmt.Exec(categoryId)
			if err != nil {
				return err
			}

		}
	}

	err = addCategory(post, postId)
	if err != nil {
		return err
	}
	return nil
}

// Categories many-to-many link
func addCategory(post Post, postId int) error {
	var postCount int

	for i := range post.Categories {
		var categoryId int64

		err := DB.QueryRow(`SELECT COUNT(*) FROM post`).Scan(&postCount)
		if err != nil {
			return err
		}

		if checkIfCategoryExist(post.Categories[i]) {
			err = DB.QueryRow("SELECT id FROM category WHERE category = ?", post.Categories[i]).Scan(&categoryId)

		} else {
			stmt, err := DB.Prepare(`INSERT INTO category (category) VALUES (?)`)
			if err != nil {
				return err
			}
			result, err := stmt.Exec(post.Categories[i])
			if err != nil {
				return err
			}

			categoryId, err = result.LastInsertId()
			if err != nil {
				return err
			}

			if postCount == 1 {
				err = insertPostCategory(postId, categoryId)
				if err != nil {
					return err
				}
			}
		}

		if postCount > 1 {
			err = insertPostCategory(postId, categoryId)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func insertPostCategory(postId int, categoryId int64) error {
	stmt, err := DB.Prepare(`INSERT INTO post_category (post_id, category_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(postId, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func SelectAllCategories() ([]Category, error) {
	var categories []Category
	rows, err := DB.Query("SELECT * FROM category")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category Category

		err = rows.Scan(&category.Id, &category.Category)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func SelectAllPostCategory() ([]PostCategory, error) {
	var post_categories []PostCategory
	rows, err := DB.Query("SELECT * FROM post_category")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post_category PostCategory

		err = rows.Scan(&post_category.PostId, &post_category.CategoryId)
		if err != nil {
			return nil, err
		}

		post_categories = append(post_categories, post_category)
	}

	return post_categories, nil
}

func deletePostLikes(postId int) error {
	stmt, err := DB.Prepare(`
		DELETE FROM like WHERE PostId = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(postId)
	if err != nil {
		return err
	}

	stmt, err = DB.Prepare(`
		DELETE FROM dislike WHERE PostId = ?
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

func checkIfPostExist(commentId int) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE id=?)", commentId).Scan(&exists)

	return err == nil && exists
}

func checkIfCategoryExist(category string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM category WHERE category=?)", category).Scan(&exists)

	return err == nil && exists
}
