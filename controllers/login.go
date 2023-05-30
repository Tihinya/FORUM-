package controllers

import (
	"forum/login"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	userId := 228 // test user id
	login.AddLogin(w, userId)
}
func LogOut(w http.ResponseWriter, r *http.Request) {}
