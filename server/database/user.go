package database

import (
	"database/sql"
	"fmt"
	"strings"
)

func CreateUser(user UserInfo) (int, error) {
	sqlStmt, err := DB.Prepare(`
	INSERT INTO users(
		email,
		username,
		password,
		profile_picture,
		role_id)
	VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}

	result, err := sqlStmt.Exec(user.Email, user.Username, user.Password, user.ProfilePicture, user.RoleID)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func SelectAllUsers() ([]UserInfo, error) {
	rows, err := DB.Query("SELECT user_id, email, username, profile_picture, role_id FROM users")
	if err != nil {
		return nil, err
	}
	var users []UserInfo
	for rows.Next() {
		var user UserInfo
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.ProfilePicture, &user.RoleID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func SelectUser(userID int) (*UserInfo, error) {
	row := DB.QueryRow("SELECT user_id, email, username, profile_picture, role_id FROM users WHERE user_id=?", userID)

	var user UserInfo
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.ProfilePicture, &user.RoleID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserPassword(userID int) (*UserInfo, error) {
	row := DB.QueryRow("SELECT password FROM users WHERE user_id=?", userID)

	var user UserInfo
	err := row.Scan(&user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(username string, email string, role int, userID int) error {
	// Generate SQL query based on provided fields
	query := "UPDATE users SET"
	var params []interface{}

	if username != "" {
		query += " username=?,"
		params = append(params, username)
	}
	if email != "" {
		query += " email=?,"
		params = append(params, email)
	}
	if role >= 0 {
		query += " role_id=?,"
		params = append(params, role)
	}

	// Remove the trailing comma
	query = strings.TrimSuffix(query, ",")

	query += " WHERE user_id=?"

	// Add the user ID to the params slice
	params = append(params, userID)

	// Execute the SQL query
	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(userID int) error {
	stmt, err := DB.Prepare(`
		DELETE FROM users
		WHERE user_id=?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}

func ReadUserLikedPosts(userID int) ([]Post, error) {
	likedPosts := make([]Post, 0)

	// Get the username for the given userID
	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	// Query the database for liked posts
	rows, err := DB.Query(`
		SELECT p.id, p.title, p.content, p.profile_picture, p.username, p.creation_date,
		p.likes, p.dislikes, p.last_edited
		FROM like AS l
		INNER JOIN post AS p ON l.post_id = p.id
		WHERE l.Username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post

		// Scan the post data from the query result
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
		if err != nil {
			return nil, err
		}

		post.Likes, _ = getPostLikes(post.Id)
		post.Dislikes, _ = getPostDislikes(post.Id)

		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)
		post.CommentCount = getCommentsCount(post.Id)
		post.UserInfo.ProfilePicture, _ = GetAvatar(post.UserInfo.Username)

		likedPosts = append(likedPosts, post)
	}

	return likedPosts, nil
}

func ReadUserDislikedPosts(userID int) ([]int, error) {
	posts := make([]int, 0)

	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT post_id, comment_id
		FROM dislike WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dislike Dislike

		err = rows.Scan(&dislike.PostId, &dislike.CommentId)
		if err != nil {
			return nil, err
		}

		if dislike.PostId != 0 {
			post := dislike.PostId
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func ReadUserLikedComments(userID int) ([]int, error) {
	posts := make([]int, 0)

	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT post_id, comment_id
		FROM like WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var like Like

		err = rows.Scan(&like.PostId, &like.CommentId)
		if err != nil {
			return nil, err
		}

		if like.CommentId != 0 {
			post := like.CommentId
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func ReadUserDislikedComments(userID int) ([]int, error) {
	posts := make([]int, 0)

	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT post_id, comment_id
		FROM dislike WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dislike Dislike

		err = rows.Scan(&dislike.PostId, &dislike.CommentId)
		if err != nil {
			return nil, err
		}

		if dislike.CommentId != 0 {
			post := dislike.CommentId
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func ReadUserCreatedPosts(userID int) ([]Post, error) {
	posts := make([]Post, 0)

	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT id, title, content, profile_picture, username, creation_date,
		likes, dislikes, last_edited
		FROM post WHERE username = ?
	`, username)
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

func ReadUserCreatedComments(userID int) ([]Comment, error) {
	comments := make([]Comment, 0)

	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT id, post_id, content, profile_picture, username, creation_date,
		likes, dislikes, last_edited, post_id
		FROM comment WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment Comment

		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.Content,
			&comment.UserInfo.ProfilePicture,
			&comment.UserInfo.Username,
			&comment.CreationDate,
			&comment.Likes,
			&comment.Dislikes,
			&comment.LastEdited,
			&comment.PostId,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func ReadUserCommentsPosts(userID int) ([]Post, error) {
	postsWithComments := make([]Post, 0)

	// Get the username for the given userID
	username, err := GetUsername(userID)
	if err != nil {
		return nil, err
	}

	// Query the database for posts with comments by the user
	rows, err := DB.Query(`
		SELECT p.id, p.title, p.content, p.profile_picture, p.username, p.creation_date,
		p.likes, p.dislikes, p.last_edited
		FROM comment AS c
		INNER JOIN post AS p ON c.post_id = p.id
		WHERE c.username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post

		// Scan the post data from the query result
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
		if err != nil {
			return nil, err
		}

		post.Likes, _ = getPostLikes(post.Id)
		post.Dislikes, _ = getPostDislikes(post.Id)

		post.Comments = fmt.Sprintf("https://localhost:8080/comments/%d", post.Id)
		post.CommentCount = getCommentsCount(post.Id)
		post.UserInfo.ProfilePicture, _ = GetAvatar(post.UserInfo.Username)

		postsWithComments = append(postsWithComments, post)
	}

	return postsWithComments, nil
}

func GetUsername(userID int) (string, error) {
	var username string

	err := DB.QueryRow("SELECT username FROM users WHERE user_id = ?", userID).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func GetAvatar(username string) (string, error) {
	var avatar string

	err := DB.QueryRow("SELECT profile_picture FROM users WHERE username = ?", username).Scan(&avatar)
	if err != nil {
		return "", err
	}

	return avatar, nil
}

func AddRole(roleName string) error {
	// Prepare the SQL statement
	stmt, err := DB.Prepare("INSERT INTO roles (name) VALUES (?)")
	if err != nil {
		return err
	}

	// Execute the SQL statement to insert the role
	_, err = stmt.Exec(roleName)
	if err != nil {
		return err
	}

	return nil
}

func GetRoleId(roleName string) (int, error) {
	row := DB.QueryRow("SELECT role_id FROM roles WHERE name = ?", roleName)

	var roleId int
	err := row.Scan(&roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Role with the specified name not found
			return 0, nil
		}
		return 0, err
	}

	return roleId, nil
}
