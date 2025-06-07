package main

import (
	"encoding/json"
	"time"
)

// Event represents the socket event
type Event struct {
	Type    string          `json:"type"`    // event type
	Payload json.RawMessage `json:"payload"` // event payload
}

// EventHandler defines a function type that handles an Event from a client.
//
// It takes two parameters:
//   - event: the Event object representing the incoming data or action.
//   - c: a pointer to the Client that sent the event.
//
// It returns an error if the event handling fails.
type EventHandler func(event Event, c *Client) error

const (
	EventSendMesasge = "send_message"
	EventNewMesasge  = "new_message"
)

// Send Message Event represents the `send_message` socket event
type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

// NewMesageEvent represents the chat event `new_message`
type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}
