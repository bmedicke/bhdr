package homeassistant

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// Config for the connection:
type Config struct {
	Scheme string `json:"scheme"`
	Server string `json:"server"`
	Token  string `json:"token"`
}

const path string = "/api/websocket"

// GetEvents published to a channel.
func GetEvents(config Config, channel chan string) {
	// TODO clean up this entire function.
	// TODO add proper error handling.
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
			"id":   1,
			"type": "subscribe_events",
		},
	)

	// listen for messages (blocking)
	// and publish them on the channel:
	for {
		m := getMessage(connection)
		channel <- m
	}
}

// synchronous message fetching:
func getMessage(connnection *websocket.Conn) string {
	message := make(map[string]interface{})
	connnection.ReadJSON(&message)
	bytestring, _ := json.Marshal(message)
	return string(bytestring)
}
