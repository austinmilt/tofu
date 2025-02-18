package model

import (
	"time"
)

type Takeover struct {
	Tier     TakeoverTier
	GameId   GameId
	PlayerId PlayerId
	Start    time.Time
	End      time.Time
}

// whether the takeover starts after now
func (t Takeover) InFuture() bool {
	return time.Now().Before(t.Start)
}

func (t Takeover) IsActive() bool {
	now := time.Now()
	return now.After(t.Start) && now.Before(t.End)
}

type TakeoverTier uint8

const (
	TakeoverTier1 TakeoverTier = iota
	TakeoverTier2
	TakeoverTier3
)
