package models

import "time"

type ReadMessageHistoryEvent struct {
	ReceiverID int `json:"receiver_id"`
}

type ReceiveMessageEvent struct {
	Message    string `json:"message"`
	ReceiverID int    `json:"receiver_id"`
}

type SendMessageEvent struct {
	ReceiveMessageEvent
	SenderID int       `json:"sender_id"`
	SentTime time.Time `json:"sent_time"`
}

type ConnectedUserListEvent struct {
	List []int `json:"list"`
}

type OrderedUserListEvent struct {
	List []int `json:"list"`
}

type ReceivedMessage struct {
	Message    string    `json:"message"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Sent       time.Time `json:"sent_time"`
}

type IsTypingEvent struct {
	SenderID     int  `json:"sender_id"`
	ReceiverID   int  `json:"receiver_id"`
	TypingStatus bool `json:"typing_status"`
}
