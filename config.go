package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/video-sharing-platform-thingy/backend/util"
)

const configFilePath = "config.json"

type googleOauthConfig struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type oauthConfig struct {
	Google googleOauthConfig `json:"google"`
}

type config struct {
	Name    string `json:"name"`
	BaseURL string `json:"baseurl"`
	Port    string `json:"port"`

	SessionCookieName string `json:"sessionCookieName"`
	LogSession        bool   `json:"logSession"`

	Oauth oauthConfig `json:"oauth"`
}

// loadConfig loads configuration data from a file.
func loadConfig() config {
	// Open the configuration file.
	file, err := os.Open("config.json")
	util.CheckError(err)
	log.Println("Opened config file")
	defer file.Close()

	// Read all of the bytes and parse the json.
	bytes, err := ioutil.ReadAll(file)
	util.CheckError(err)
	var loadedConfig config
	json.Unmarshal(bytes, &loadedConfig)

	return loadedConfig
}
