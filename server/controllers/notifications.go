package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"forum/database"
	"forum/session"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := session.GetUserId(r)
	if err != nil {
		ReturnMessageJSON(w, "Error fetching user id", http.StatusInternalServerError, "error")
		return
	}

	notifications, err := database.GetNotifications(userId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	json.NewEncoder(w).Encode(notifications)
}

func MarkNotificationsRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userId, err := session.GetUserId(r)
	if err != nil {
		ReturnMessageJSON(w, "Error fetching user id", http.StatusInternalServerError, "error")
		return
	}

	err = database.MarkNotificationRead(userId)
	if err != nil {
		log.Println(err)
		ReturnMessageJSON(w, "Internal server error", http.StatusInternalServerError, "error")
		return
	}

	ReturnMessageJSON(w, "Notifications marked as read", http.StatusOK, "success")
}
