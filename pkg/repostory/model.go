package repository

import (
	"time"

	"github.com/google/uuid"
)

const (
	AuthEvent    = "auth"
	ErrorEvent   = "error"
	MessageEvent = "message"
	TypingEvent  = "typing"
)

type Channel uuid.UUID

type EventType string

type Sender struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
}

type Message struct {
	ID     uuid.UUID   `json:"id"`
	SendAt time.Time   `json:"send_at"`
	Event  EventType   `json:"event"`
	Sender Sender      `json:"sender"`
	Data   interface{} `json:"data"`
}

func ErrorMessage(err error) *Message {
	return &Message{
		ID:    uuid.New(),
		Event: ErrorEvent,
		Data:  err.Error(),
	}
}

func TextMessage(msg string, sender Sender) *Message {
	return &Message{
		ID:     uuid.New(),
		Event:  MessageEvent,
		Sender: sender,
		Data:   msg,
	}
}
