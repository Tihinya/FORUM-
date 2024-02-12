package socket_test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"forum/controllers"
	"forum/database"
	"forum/models"
	"forum/session"
	"forum/socket"

	"github.com/gorilla/websocket"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    int    `json:"from"`
	To      int    `json:"to"`
}

// Testing accounts follow format username@example.com:username

const databasePath = "mock_database.db"

func TestSecureWebSocketConnection(t *testing.T) {
	brothermanAcc := []string{"brotherman@example.com", "brotherman"}
	session.NewStorage()

	server := httptest.NewTLSServer(http.HandlerFunc(socket.NewManager().ServeWS))
	defer server.Close()

	// Convert http://127.0.0.1 to wss://127.
	url := "wss" + server.URL[5:]

	user := connectSocket(t, brothermanAcc, url)
	defer user.Close()
}

// TestWebSocketConversation tests sending and receiving messages over WebSocket.
func TestWebSocketConversation(t *testing.T) {
	brothermanAcc := []string{"brotherman@example.com", "brotherman"}
	testAcc := []string{"test@example.com", "test"}
	abobisAcc := []string{"abobis@example.com", "abobis"}
	session.NewStorage()

	server := httptest.NewTLSServer(http.HandlerFunc(socket.NewManager().ServeWS))
	defer server.Close()

	url := "wss" + server.URL[5:]

	// First user connection
	user1 := connectSocket(t, brothermanAcc, url)
	defer user1.Close()
	user1username := strings.Split(brothermanAcc[0], "@")[0]
	user1id, _ := database.GetUserId(user1username)

	// Second user connection
	user2 := connectSocket(t, testAcc, url)
	defer user2.Close()
	user2username := strings.Split(testAcc[0], "@")[0]
	user2id, _ := database.GetUserId(user2username)

	// Third user connection
	user3 := connectSocket(t, abobisAcc, url)
	defer user2.Close()
	user3username := strings.Split(abobisAcc[0], "@")[0]
	user3id, _ := database.GetUserId(user3username)

	message := "Hello World!"
	sendMessage(t, user1, message, user1id, user2id) // send message from user id to user id

	readMessage(t, user2, message, user2id, 0, nil) // check if given user id received message meant for him
	readMessage(t, user3, message, user3id, 0, nil) // check if given user id received message not meant for him
}

func TestWebSocketOnlineList(t *testing.T) {
	abobisAcc := []string{"abobis@example.com", "abobis"}
	session.NewStorage()

	server := httptest.NewTLSServer(http.HandlerFunc(socket.NewManager().ServeWS))
	defer server.Close()

	url := "wss" + server.URL[5:]

	user := connectSocket(t, abobisAcc, url)
	defer user.Close()

	userusername := strings.Split(abobisAcc[0], "@")[0]
	userid, _ := database.GetUserId(userusername)

	readMessage(t, user, "", 0, userid, nil)
}

func connectSocket(t *testing.T, account []string, url string) *websocket.Conn {
	// Create a dialer with TLS certificate verification disabled
	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// First user cookie
	cookie := loginAndGetCookie(t, account)
	header := http.Header{}
	header.Add("Cookie", "session_token="+cookie.Value)

	// Connect user
	user1, _, err := dialer.Dial(url, header)
	if err != nil {
		t.Fatalf("Could not open a ws connection for user: %v", err)
	}

	username := strings.Split(account[0], "@")[0]
	userid, _ := database.GetUserId(username)
	fmt.Printf("Successfully connected account %v with id %v\n", account[0], userid)

	return user1
}

// Login & get session_token & open database.db for other database related functions to work
func loginAndGetCookie(t *testing.T, account []string) *http.Cookie {
	// Open database file from given path and set database.DB to it
	err := database.OpenDatabase(databasePath)
	if err != nil {
		t.Fatal("Failed opening db", err)
	}

	// Replace with your credentials
	credentials := map[string]string{
		"email":    account[0],
		"password": account[1],
	}
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request to simulate an HTTP request
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the Login function directly
	controllers.Login(w, req)

	// Convert the result into an http.Response object
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed: %s", resp.Status)
	}

	// Extract the cookie from the response
	var sessionTokenCookie *http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "session_token" {
			sessionTokenCookie = c
			break
		}
	}

	if sessionTokenCookie == nil {
		t.Fatal("Cookie not found in the response")
	}

	if sessionTokenCookie.Value == "" {
		t.Fatal("Session token empty")
	}

	return sessionTokenCookie
}

func sendMessage(t *testing.T, connection *websocket.Conn, message string, senderId int, recipientId int) {
	sendData := &models.SendMessageEvent{
		ReceiveMessageEvent: models.ReceiveMessageEvent{
			Message:    message,
			ReceiverID: recipientId,
		},
		SenderID: senderId,
		SentTime: time.Now(),
	}

	jsonData, err := json.Marshal(sendData)
	if err != nil {
		t.Fatal("Failed marshalling to JSON", err)
	}

	var payload json.RawMessage = jsonData

	eventData := socket.Event{
		Type:    "send_message",
		Payload: payload,
	}

	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		t.Fatal("Failed marshalling to JSON", err)
	}

	err = connection.WriteMessage(websocket.TextMessage, eventJSON)
	if err != nil {
		t.Fatalf("Could not send message from user 1 to user 2 %v", err)
	}
	fmt.Printf("Sent message: %v to id:%v from id:%v\n", message, recipientId, senderId)
}

func readMessage(t *testing.T, connection *websocket.Conn, sentMessage string, receiverId int, onlineUsersListId int, msgHistoryFromDB []models.SendMessageEvent) {
	timeout := time.After(1 * time.Second)
	messageChan := make(chan []byte)
	errChan := make(chan error)

	// Start a goroutine to read messages
	go func() {
		for {
			_, msg, err := connection.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			messageChan <- msg
		}
	}()

	for {
		select {
		case <-timeout:
			fmt.Printf("Timeout: no expected message received within 1 second for user id %v\n", receiverId)
			return

		case err := <-errChan:
			t.Fatalf("Could not read message: %v", err)

		case receivedPayload := <-messageChan:
			var event socket.Event
			if err := json.Unmarshal(receivedPayload, &event); err != nil {
				t.Fatalf("bad payload in request: %v", err)
			}

			switch event.Type {
			case socket.EventReceiveMessage:
				if sentMessage != "" {
					var message models.SendMessageEvent
					if err := json.Unmarshal(event.Payload, &message); err != nil {
						t.Fatalf("Error unmarshaling message: %v\n", err)
					}
					if sentMessage == message.Message {
						fmt.Printf("Received message %v to id %v\n", message.Message, message.SenderID)
						return
					}
				}

			case socket.EventOnlineUserList:
				if onlineUsersListId != 0 {
					var onlineUsersList models.ConnectedUserListEvent
					if err := json.Unmarshal(event.Payload, &onlineUsersList); err != nil {
						t.Fatalf("Error unmarshaling online users list: %v\n", err)
					}
					if !contains(onlineUsersList.List, onlineUsersListId) {
						t.Fatalf("Expected user id %v to be in online user list: %v\n", onlineUsersListId, onlineUsersList.List)
					}
					fmt.Printf("Found user id %v in online users list %v\n", onlineUsersListId, onlineUsersList.List)
					return
				}

			case socket.EventReceiveMessageHistory:
				var messageHistory []models.SendMessageEvent
				if err := json.Unmarshal(event.Payload, &messageHistory); err != nil {
					t.Fatalf("Error unmarshaling message history list: %v\n", err)
				}

				if !reflect.DeepEqual(messageHistory, msgHistoryFromDB) {
					t.Fatalf("Arrays are not equal: %v", event.Type)
				}
				fmt.Println("Message history arrays are equal")
				return

			default:
				t.Fatalf("Unexpected message type: %v", event.Type)
			}
		}
	}
}

func contains(arr []int, userid int) bool {
	for _, id := range arr {
		if userid == id {
			return true
		}
	}
	return false
}
