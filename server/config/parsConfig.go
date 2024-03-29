package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var Config = parseConfig()

type config struct {
	Port               string `json:"port"`
	ProfilePicture     string `json:"profile_picture"`
	GitHubClientId     string `json:"gitHubClientId"`
	GitHubClientSecret string `json:"gitHubClientSecret"`
	GitHubRedirectURI  string `json:"gitHubRedirectURI"`
	GoogleID           string `json:"googleID"`
	GoogleClientSecret string `json:"googleClientSecret"`
	GoogleRedirectURI  string `json:"googleRedirectURI"`
	GoogleOAuth        string `json:"googleOAuth"`
	GoogleGetToken     string `json:"googleGetToken"`
}

func parseConfig() config {
	var config config

	jsonFile, err := os.Open("dev_config.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(byteValue, &config)
	return config
}
