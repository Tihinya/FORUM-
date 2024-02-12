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

	err = postAddCategory(post, int(postId))
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

func CreatePostReport(post PostReport, userId int) error {
	stmt, err := DB.Prepare(`
	INSERT INTO post_reports (message, status, seen, user_id, post_id)
    VALUES (?, 'pending', false, ?, ?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Message, userId, post.PostID)
	if err != nil {
		return err
	}

	return nil
}
func SelectAllPostReports() ([]PostReport, error) {
	var reports []PostReport

	// Query to select all reports from the "post_reports" table along with the "title" from the "post" table and the username from the "users" table
	query := `
		SELECT r.post_id, r.report_id, r.message, r.status, p.title, u.username
		FROM post_reports r
		INNER JOIN post p ON r.post_id = p.id
		INNER JOIN users u ON r.user_id = u.user_id
		WHERE r.status = 'pending'
	`

	// Execute the query
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan data into PostReport struct
	for rows.Next() {
		var report PostReport
		if err := rows.Scan(&report.PostID, &report.ReportID, &report.Message, &report.Status, &report.Title, &report.UserName); err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	// Check for errors in row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func UpdatePostReport(report PostReport) error {
	// Query to update a report's content and status
	query := "UPDATE post_reports SET response = ?, status = ? WHERE report_id = ?"

	// Execute the update query
	_, err := DB.Exec(query, report.Response, report.Status, report.ReportID)
	if err != nil {
		return err
	}

	return nil
}
func UpdatePostReportAnswer(report PostReport) error {
	// Query to update a report's content and status
	query := "UPDATE post_reports SET seen = ? WHERE report_id = ?"

	// Execute the update query
	_, err := DB.Exec(query, report.Seen, report.ReportID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserReportedPost(UserID int) ([]PostReport, error) {
	// Initialize an empty slice to store the reported posts
	var reportedPosts []PostReport

	// Query to select reported posts by a specific user
	query := `
        SELECT r.post_id, r.report_id, r.message, r.status, p.title
        FROM post_reports r
        INNER JOIN post p ON r.post_id = p.id
        WHERE r.user_id = ? AND r.status IN (?, ?) AND r.seen = ?;
    `

	// Execute the query
	rows, err := DB.Query(query, UserID, "approved", "rejected", 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan data into PostReport struct
	for rows.Next() {
		var report PostReport
		if err := rows.Scan(&report.PostID, &report.ReportID, &report.Message, &report.Status, &report.Title); err != nil {
			return nil, err
		}

		reportedPosts = append(reportedPosts, report)
	}

	// Check for errors in row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reportedPosts, nil
}
