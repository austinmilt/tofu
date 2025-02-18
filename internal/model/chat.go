package model

import "time"

type ChatEvent struct {
	Chat  *ChatMessage
	Clear bool
}

type ChatMessage struct {
	ChatterHandle string
	Message       string
	Timestamp     time.Time
}

type EventConsumerId = string
