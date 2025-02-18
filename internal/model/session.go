package model

import "time"

type SessionId = string

type PlayerSessionRecord struct {
	// keys are the raw validated input strings sent by the player
	InputCounts map[string]uint32

	// number of votes which were part of a consensus and thus were
	// sent to the game. Note that here vote is the same as input
	// in InputCounts
	ConsensusVoteCount uint32
	Start              time.Time
	End                time.Time
}
