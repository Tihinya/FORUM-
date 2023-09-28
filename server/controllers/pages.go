package controllers

import (
	"fmt"
	"net/http"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Hello, World!</h1></body></html>")
}
func LoginRegistrationPage(w http.ResponseWriter, r *http.Request) {}
func ProfilePage(w http.ResponseWriter, r *http.Request)           {}
func ErrorPage(w http.ResponseWriter, r *http.Request)             {}
func PerformancePage(w http.ResponseWriter, r *http.Request)       {}
