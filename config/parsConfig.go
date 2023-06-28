package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var Config = parseConfig()

type config struct {
	Port           string `json:"port"`
	ProfilePicture string `json:"profile_picture"`
}

func parseConfig() config {
	var config config

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
