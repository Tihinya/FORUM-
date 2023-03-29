package controllers

import (
	"forum/login"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := login.AddLogin(w, 228)
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func LogOut(w http.ResponseWriter, r *http.Request) {}
