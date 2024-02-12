package socket_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/database"
	"forum/models"
	"forum/socket"
	"forum/validation"

	"github.com/gorilla/websocket"
)

func TestReadMessagesHandler(t *testing.T) {
	adminCredentials := []string{"admin@example.com", "admin"}

	testServer := httptest.NewTLSServer(http.HandlerFunc(socket.NewManager().ServeWS))
	defer testServer.Close()

	websocketURL := "wss" + testServer.URL[5:]
	socketConnection := connectSocket(t, adminCredentials, websocketURL)

	OtinsID, AdminID := getUserIDs(t, "abobis@example.com", "admin@example.com")

	defer func() {
		deleteMessages(t, OtinsID, AdminID)
		deleteMessages(t, AdminID, OtinsID)
	}()

	createMessage(t, "Admin test message 1", OtinsID, AdminID)
	createMessage(t, "Otis test message 2", AdminID, OtinsID)
	createMessage(t, "Admin test message 3", OtinsID, AdminID)

	msgHistory := readMessageHistoryFromDB(AdminID, OtinsID, t)

	sendMessagesHistoryRequest(t, socketConnection, OtinsID)

	readMessage(t, socketConnection, "", 0, 0, msgHistory)
}

func getUserIDs(t *testing.T, otisEmail, adminEmail string) (int, int) {
	OtinsID, err := validation.GetUserID(database.DB, otisEmail, "")
	if err != nil {
		t.Fatalf("Failed to get Otis ID: %v", err)
	}

	AdminID, err := validation.GetUserID(database.DB, adminEmail, "")
	if err != nil {
		t.Fatalf("Failed to get Admin ID: %v", err)
	}

	return OtinsID, AdminID
}

func deleteMessages(t *testing.T, senderID, recipientID int) {
	err := database.DeleteMessages(senderID, recipientID)
	if err != nil {
		t.Fatalf("Failed to delete messages from database: %v", err)
	}
}

func createMessage(t *testing.T, message string, senderID, recipientID int) {
	_, err := database.CreateMessage(message, senderID, recipientID)
	if err != nil {
		t.Fatalf("Failed to write message into database: %v", err)
	}
}

func sendMessagesHistoryRequest(t *testing.T, connection *websocket.Conn, from int) {
	testReadMessageEvent := models.ReadMessageHistoryEvent{
		ReceiverID: from,
	}

	payload, err := json.Marshal(testReadMessageEvent)
	if err != nil {
		t.Fatalf("Failed to marshal ReadMessageEvent payload: %v", err)
	}

	event := socket.Event{
		Type:    socket.EventReadMessagesHistory,
		Payload: payload,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		t.Fatal("Failed marshalling to JSON", err)
	}

	err = connection.WriteMessage(websocket.TextMessage, eventJSON)
	if err != nil {
		t.Fatalf("Could not send message %v", err)
	}
}

func readMessageHistoryFromDB(sender, recipient int, t *testing.T) []models.SendMessageEvent {
	messages, err := database.ReadMessage(recipient, sender)
	if err != nil {
		t.Fatalf("Failed to read response from WebSocket server: %v", err)
	}

	return messages
}
