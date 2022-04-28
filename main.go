package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bmedicke/bhdr/util"
)

//go:embed bhdr.json
var bhdrJSON string

func main() {
	// register and parse flags:
	createConfig := flag.Bool(
		"create-config",
		false,
		"create bhdr.json config file in $HOME",
	)
	flag.Parse()

	// get user's home folder:
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("home folder error: ", err)
	}
	configFile := filepath.Join(home, "bhdr.json")

	// handle --create-config flag:
	if *createConfig {
		err := util.CreateFileIfNotExist(configFile, bhdrJSON)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("file %v created\n", configFile)
		os.Exit(0)
	}

	// read config file:
	haConfigFile, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err, ". you can create one with: bhdr --create-config")
	}

	// unmarshal config:
	var config map[string]interface{}
	jsonParser := json.NewDecoder(haConfigFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal("config file parsing error: ", err)
	}

	spawnTUI(config)
}
