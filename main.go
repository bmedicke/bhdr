package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/bmedicke/bhdr/homeassistant"
)

func main() {
	// get user's home folder:
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("home folder error: ", err)
	}

	// read config file:
	haConfigFile, err := os.Open(filepath.Join(home, "bhdr.json"))
	if err != nil {
		log.Fatal(err, ". you can create one with: bhdr init")
	}

	// unmarshal config:
	var haConfig homeassistant.Config
	jsonParser := json.NewDecoder(haConfigFile)
	err = jsonParser.Decode(&haConfig)
	if err != nil {
		log.Fatal("config file parsing error: ", err)
	}

	spawnTUI(haConfig)
}
