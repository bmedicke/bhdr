package homeassistant

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// Connect connects to Home Assistant and communicates with two channels:
// * events: events from HA will be published here
// * commands: commands will be sent to HA
func Connect(config Config, events chan string, commands chan Command) {
	// TODO add proper error handling.
	var messageID uint = 1
	const APIPath string = "/api/websocket"
	haURL := url.URL{
		Scheme: config.Scheme,
		Host:   config.Server,
		Path:   APIPath,
	}

	// connect:
	connection, _, err := websocket.DefaultDialer.Dial(haURL.String(), nil)
	defer connection.Close()
	if err != nil {
		log.Fatal(err)
	}

	// authenticate:
	connection.WriteJSON(
		map[string]string{
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
		command := <-commands
		haCommand := map[string]interface{}{}

		if command.EntityID != "" {
			haCommand["target"] = map[string]string{
				"entity_id": command.EntityID,
			}
		}

		if command.Service != "" {
			haCommand["service"] = command.Service
		}

		if command.Domain {
			haCommand["domain"] = strings.Split(command.EntityID, ".")[0]
		}

		haCommand["type"] = command.Type
		haCommand["id"] = messageID

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
