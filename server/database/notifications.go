package database

import (
	"time"
)

func GetNotifications(userId int) ([]Notification, error) {
	notifications := make([]Notification, 0)

	username, err := GetUsername(userId)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
		SELECT * FROM notifications WHERE username = ?
	`, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var notification Notification

		err = rows.Scan(
			&notification.Id,
			&notification.Username,
			&notification.RelatedObjectType,
			&notification.RelatedObjectId,
			&notification.Type,
			&notification.Status,
			&notification.CreationDate,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func CreateNotification(objectType string, objectId int, Type string) error {
	username, err := getPostCreatorByPostId(objectId)
	if err != nil {
		return err
	}

	stmt, err := DB.Prepare(`
		INSERT INTO notifications (
			username,
			related_object_type,
			related_object_id,
			type,
			creation_date
		) VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username, objectType, objectId, Type, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func MarkNotificationRead(userId int) error {
	username, err := GetUsername(userId)
	if err != nil {
		return err
	}

	stmt, err := DB.Prepare(`
		UPDATE notifications SET
			status = ?
		WHERE username = ?
	`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec("read", username)
	if err != nil {
		return err
	}

	return nil
}
