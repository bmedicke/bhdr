package homeassistant

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// Config for the connection:
type Config struct {
	Scheme string `json:"scheme"`
	Server string `json:"server"`
	Token  string `json:"token"`
}

const path string = "/api/websocket"

// Connect connects to Home Assistant and communicates with two channels:
// * events: events from HA will be published here
// * commands: commands will be sent to HA
func Connect(config Config, events chan string, commands chan string) {
	// TODO clean up this entire function.
	// TODO add proper error handling.
	var messageID uint = 1
	haURL := url.URL{Scheme: config.Scheme, Host: config.Server, Path: path}

	// connect:
	connection, _, err := websocket.DefaultDialer.Dial(haURL.String(), nil)
	defer connection.Close()
	if err != nil {
		log.Fatal(err)
	}

	// authenticate:
	connection.WriteJSON(
		map[string]interface{}{
			"type":         "auth",
			"access_token": config.Token,
		},
	)

	// subscribe to all:
	connection.WriteJSON(
		map[string]interface{}{
			"id":   messageID,
			"type": "subscribe_events",
		},
	)
	messageID++

	// listen for messages from HA and publish them on the events channel:
	go func(events chan string, connection *websocket.Conn) {
		for {
			events <- getMessage(connection)
		}
	}(events, connection)

	// listen for commands and send them to HA:
	for {
		entityID := <-commands
		target := map[string]interface{}{
			"entity_id": entityID,
		}
		domain := strings.Split(entityID, ".")[0]

		haCommand := map[string]interface{}{
			"id":      messageID,
			"type":    "call_service",
			"domain":  domain,
			"service": "toggle",
			"target":  target,
		}

		connection.WriteJSON(haCommand)
		messageID++
	}
}

// synchronous message fetching:
func getMessage(connnection *websocket.Conn) string {
	message := make(map[string]interface{})
	connnection.ReadJSON(&message)
	bytestring, _ := json.Marshal(message)

	var pretty bytes.Buffer
	json.Indent(&pretty, bytestring, "", "  ")
	return string(pretty.Bytes())
}
