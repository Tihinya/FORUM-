package database

func getCategories(post Post) ([]string, error) {
	categories := make([]string, 0)

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

		categories = append(categories, category)
	}

	return categories, nil
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
			if err != nil {
				return err
			}
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
	categories := make([]Category, 0)

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
	postCategories := make([]PostCategory, 0)

	rows, err := DB.Query("SELECT * FROM post_category")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var postCategory PostCategory

		err = rows.Scan(&postCategory.PostId, &postCategory.CategoryId)
		if err != nil {
			return nil, err
		}

		postCategories = append(postCategories, postCategory)
	}

	return postCategories, nil
}

func checkIfCategoryExist(category string) bool {
	var exists bool

	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM category WHERE category = ?)", category).Scan(&exists)

	return err == nil && exists
}
