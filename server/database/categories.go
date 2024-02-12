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

// Categories many-to-many link
func postAddCategory(post Post, postId int) error {
	var postCount int

	for i := range post.Categories {
		var categoryId int64

		err := DB.QueryRow(`SELECT COUNT(*) FROM post`).Scan(&postCount)
		if err != nil {
			return err
		}

		err = DB.QueryRow("SELECT id FROM category WHERE category = ?", post.Categories[i]).Scan(&categoryId)
		if err != nil {
			return err
		}

		err = insertPostCategory(postId, categoryId)
		if err != nil {
			return err
		}

	}

	return nil
}

func AddCategory(category string) error {
	stmt, err := DB.Prepare(`INSERT INTO category (category) VALUES (?)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(category)
	if err != nil {
		return err
	}

	return nil
}

func DeleteCategory(categoryId int) error {

	// Delete all post_category rows with given category id
	stmt, err := DB.Prepare(`
		DELETE FROM post_category WHERE category_id = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(categoryId)
	if err != nil {
		return err
	}

	// Delete category row from category table with given category id
	stmt, err = DB.Prepare(`
		DELETE FROM category WHERE id = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(categoryId)
	if err != nil {
		return err
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
