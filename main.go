package main

import (
	"encoding/json"
	"io/ioutil"
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
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValue, &config)
	return config
}

func main() {

	config := ParseConfig()

	log.Println("Ctrl + Click on the link: https://localhost:" + config.Port)
	log.Println("To stop the server press `Ctrl + C`")
	log.Fatal(http.ListenAndServeTLS(":"+config.Port, "cert.pem", "key.pem", nil))
	/*
		This is a proposed endpoint design that groups actions by the data type they operate on
		because real-time-forum and social-network projects are going to be API based
		It's important to note that this is just a design and the actual implementation may vary.

		Logging, Notification, Access Level, Status, Security, Performance will be implemented as middleware on top
		of the enpoints.
	*/

	// Post:
	// --Create
	// --Read
	// --Update
	// --Delete

	// User:
	// --Create
	// --Read
	// --Update
	// --Delete
	// --Change permissions

	// Login:
	// --Login
	// --Logout
	// ( register is User:Create )

	// Comment:
	// --Create
	// --Read
	// --Update
	// --Delete

	// Pages:
	// --Main
	// --Loggin/registration
	// --Profile
	// --Error page
	// --Performance limit page
}
