package socket

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"forum/models"
	"forum/session"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ErrEventNotSupported = errors.New("this event type is not supported")
	Instance             *Manager
)

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()

	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventReadMessagesHistory] = ReadMessagesHistoryHandler
	m.handlers[EventIsTyping] = TypingStatusHandler
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Manager) getClient(userId int) *Client {
	for client := range m.clients {
		if client.id == userId {
			return client
		}
	}
	return nil
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func (m *Manager) RemoveClientWithId(id int) {
	client := m.getClient(id)
	m.removeClient(client)
	broadcastOnlineUserList(m)
}

func (m *Manager) GetConnectedClients() Event {
	var onlineUserList models.ConnectedUserListEvent

	for client := range m.clients {
		onlineUserList.List = append(onlineUserList.List, client.id)
	}

	data, err := json.Marshal(onlineUserList)
	if err != nil {
		log.Printf("failed to marshal online user list: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventOnlineUserList

	return outgoingEvent
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	websocketUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	id, err := session.GetUserId(r)
	if err != nil {
		log.Println(err)
		return
	}

	// Create New Client
	client := NewClient(conn, m, id)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()

	broadcastOnlineUserList(m)
	broadcastUsersByLastMessage(m, id)
}
