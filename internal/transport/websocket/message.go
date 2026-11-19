package websocket

import (
	"encoding/json"
)

type IncomingMessage struct {
	Type       string          `json:"type"`
	ActionType string          `json:"actionType"`
	Payload    json.RawMessage `json:"payload"`
	Raw        json.RawMessage `json:"-"`
}

type OutgoingMessage struct {
	Type    string `json:"type"`
	Entity  string `json:"entity"`
	Payload any    `json:"payload"`
	Raw     any    `json:"-"`
}

type WelcomeMessage struct {
	Type     string `json:"type"`
	PlayerId string `json:"player_id"`
	Seat     int16  `json:"seat"`
}
