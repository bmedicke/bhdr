package homeassistant

// Config for the connection:
type Config struct {
	Scheme string `json:"scheme"`
	Server string `json:"server"`
	Token  string `json:"token"`
}

// Command that can be sent to the commands channel.
type Command struct {
	EntityID string
	Service  string
	Type     string
	Domain   bool
}

// Message is the top level JSON object of a HA WS response.
type Message struct {
	Type   string   `json:"type"`
	Result []Result `json:"result"`
	Event  Event    `json:"event"`
}

// State is attached to Result and Data.
type State struct {
	State string `json:"state"`
}

// Result is an optional JSON object for Message.
type Result struct {
	State    string `json:"state"`
	EntityID string `json:"entity_id"`
}

// Event is an optional JSON object for Message.
type Event struct {
	Type string `json:"event_type"`
	Data Data   `json:"data"`
}

// Data is part of Event responses.
type Data struct {
	EntityID string `json:"entity_id"`
	NewState State  `json:"new_state"`
	OldState State  `json:"old_state"`
	NickName string
}
