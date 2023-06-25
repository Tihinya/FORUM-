package main

import (
	"encoding/json"
	"fmt"
	ct "forum/controllers"
	"forum/database"
	"forum/router"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Port string `json:"port"`
}

func ParseConfig() Config {
	var config Config

	jsonFile, err := os.Open("dev_config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, &config)
	return config
}

func ExampleMiddleware() router.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Your middleware logic here
			fmt.Println("Example middleware executed 1")

			// Call the next middleware/handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	config := ParseConfig()

	r := router.NewRouter()

	database.CreateTables()

	// middleware usage example
	r.AddGlobalMiddleware(ExampleMiddleware())

	// User
	r.NewRoute("POST", `/user/create`, ct.CreateUser)
	r.NewRoute("GET", `/user/(?P<id>\d+)/get`, ct.ReadUser)
	r.NewRoute("GET", `/users/get`, ct.ReadUsers)
	r.NewRoute("PATCH", `/user/(?P<id>\d+)/update`, ct.UpdateUser)
	r.NewRoute("DELETE", `/user/(?P<id>\d+)/delete`, ct.DeleteUser)

	// Post
	r.NewRoute("POST", `/post`, ct.CreatePost)
	r.NewRoute("GET", `/post/(?P<id>\d+)`, ct.ReadPost)
	r.NewRoute("PATCH", `/post/(?P<id>\d+)`, ct.UpdatePost)
	r.NewRoute("DELETE", `/post/(?P<id>\d+)`, ct.DeletePost)
	r.NewRoute("GET", `/posts`, ct.ReadPosts)
	r.NewRoute("GET", `/categories`, ct.ReadCategories)
	r.NewRoute("GET", `/postcategories`, ct.ReadPostCategories) // TO BE REMOVED

	// Comment
	r.NewRoute("POST", `/comment/(?P<id>\d+)`, ct.CreateComment)
	r.NewRoute("GET", `/comment/(?P<id>\d+)`, ct.ReadComment)
	r.NewRoute("PATCH", `/comment/(?P<id>\d+)`, ct.UpdateComment)
	r.NewRoute("DELETE", `/comment/(?P<id>\d+)`, ct.DeleteComment)
	r.NewRoute("GET", `/comments/(?P<id>\d+)`, ct.ReadComments)

	// Login
	r.NewRoute("GET", `/login`, ct.Login)
	r.NewRoute("GET", `/logout/(?P<id>\d+)`, ct.LogOut)

	http.HandleFunc("/", r.Serve)

	log.Println("Ctrl + Click on the link: https://localhost:" + config.Port)
	log.Println("To stop the server press `Ctrl + C`")
	log.Fatal(http.ListenAndServeTLS(":"+config.Port, "cert.pem", "key.pem", nil))
}
