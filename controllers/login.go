package controllers

import (
	"forum/login"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	login.Loginadd(w, 12)
}
func LogOut(w http.ResponseWriter, r *http.Request) {}
