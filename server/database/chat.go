package database

import (
	"log"
	"time"

	"forum/models"
)

func CreateMessage(message string, receiver int, sender int) (*models.SendMessageEvent, error) {
	creationDate := time.Now().UTC()
	_, err := DB.Exec(`INSERT INTO chat_messages (message, created_at, receiver_id, sender_id)
	VALUES (?, ?, ?, ?);`, message, creationDate, receiver, sender)
	if err != nil {
		return nil, err
	}

	return &models.SendMessageEvent{
		ReceiveMessageEvent: models.ReceiveMessageEvent{
			Message:    message,
			ReceiverID: receiver,
		},
		SenderID: sender,
		SentTime: creationDate,
	}, err
}

func ReadMessage(receiver, sender int) ([]models.SendMessageEvent, error) {
	rows, err := DB.Query(`
	SELECT message, sender_id, receiver_id, created_at
	FROM chat_messages
	WHERE (sender_id = $1 AND receiver_id = $2)
	OR (sender_id = $2 AND receiver_id = $1)
	ORDER BY created_at
`, sender, receiver)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]models.SendMessageEvent, 0)
	for rows.Next() {
		var message models.SendMessageEvent

		err := rows.Scan(&message.Message, &message.SenderID, &message.ReceiverID, &message.SentTime)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func OrderUserIdsByLastMessage(sender int) ([]int, error) {
	rows, err := DB.Query(`
	SELECT message, sender_id, receiver_id, created_at
	FROM chat_messages
	WHERE sender_id = $1
	OR receiver_id = $1
	ORDER BY created_at DESC;
	`, sender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderedUserIds := make([]int, 0)
	for rows.Next() {
		var message models.SendMessageEvent
		var found bool

		err := rows.Scan(&message.Message, &message.SenderID, &message.ReceiverID, &message.SentTime)
		if err != nil {
			return nil, err
		}

		for _, userId := range orderedUserIds {
			if message.SenderID == userId || message.ReceiverID == userId {
				found = true
			}
		}

		if found || message.SenderID == 0 || message.ReceiverID == 0 {
			continue
		}

		if message.SenderID == sender {
			orderedUserIds = append(orderedUserIds, message.ReceiverID)
		} else {
			orderedUserIds = append(orderedUserIds, message.SenderID)
		}
	}

	return orderedUserIds, nil
}

func DeleteMessages(senderID, recipientID int) error {
	stmt, err := DB.Prepare("DELETE FROM chat_messages WHERE sender_id = ? AND receiver_id = ?")
	if err != nil {
		log.Println("Error preparing DELETE statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(senderID, recipientID)
	if err != nil {
		log.Println("Error executing DELETE statement:", err)
		return err
	}

	return nil
}
