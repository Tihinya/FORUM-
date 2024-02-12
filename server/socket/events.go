package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"forum/database"
	"forum/models"
)

const (
	EventSendMessage               = "send_message"
	EventOnlineUserList            = "online_users_list"
	EventReceiveMessage            = "receive_message"
	EventReadMessagesHistory       = "read_messages_history"
	EventReceiveMessageHistory     = "receive_messages_history"
	EventReceiveUsersByLastMessage = "receive_users_by_last_message"
	EventIsTyping                  = "typing_status"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

func SendMessageHandler(event Event, c *Client) error {
	var chatEvent models.ReceiveMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	chatEvent.Message = strings.TrimSpace(chatEvent.Message)
	if chatEvent.Message == "" {
		return nil
	}

	sender := c.id
	recipient := chatEvent.ReceiverID

	if recipient == sender {
		return fmt.Errorf("error: failed to send messages to yourself")
	}

	message, err := database.CreateMessage(chatEvent.Message, recipient, sender)
	if err != nil {
		return fmt.Errorf("failed to write the message to the database: %v", err)
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event = Event{
		Type:    EventReceiveMessage,
		Payload: data,
	}

	for client := range c.manager.clients {
		if client.id == recipient || client.id == sender {
			client.egress <- outgoingEvent
		}
	}

	broadcastUsersByLastMessage(c.manager, sender)
	broadcastUsersByLastMessage(c.manager, recipient)

	return nil
}

func ReadMessagesHistoryHandler(event Event, c *Client) error {
	var chatEvent models.ReadMessageHistoryEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}
	sender := c.id
	recipient := chatEvent.ReceiverID

	if recipient == sender {
		return fmt.Errorf("error: failed to send messages to yourself")
	}

	messages, err := database.ReadMessage(recipient, sender)
	if err != nil {
		log.Println(event, err)
	}

	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event = Event{
		Type:    EventReceiveMessageHistory,
		Payload: data,
	}

	for client := range c.manager.clients {
		if client.id == sender {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func TypingStatusHandler(event Event, c *Client) error {
	var typingStatusEvent models.IsTypingEvent
	if err := json.Unmarshal(event.Payload, &typingStatusEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	sender := c.id
	recipient := typingStatusEvent.ReceiverID

	if recipient == sender {
		return fmt.Errorf("error: failed to send messages to yourself")
	}

	for client := range c.manager.clients {
		if client.id == recipient {
			client.egress <- event
		}
	}

	return nil
}

func broadcastOnlineUserList(m *Manager) {
	onlineUsersListEvent := m.GetConnectedClients()
	for client := range m.clients {
		client.egress <- onlineUsersListEvent
	}
}

func broadcastUsersByLastMessage(m *Manager, clientId int) {
	userIds, err := database.OrderUserIdsByLastMessage(clientId)
	if err != nil {
		fmt.Printf("failed to get ordered users by last mesage: %v\n", err)
	}
	var orderedUserIds models.OrderedUserListEvent
	orderedUserIds.List = userIds

	data, err := json.Marshal(orderedUserIds)
	if err != nil {
		fmt.Printf("failed to marshal broadcast message: %v\n", err)
	}

	var outgoingEvent Event = Event{
		Type:    EventReceiveUsersByLastMessage,
		Payload: data,
	}

	client := m.getClient(clientId)
	if client != nil {
		client.egress <- outgoingEvent
	}
}
