package model

import "time"

type Report struct {
	ReportedPlayerId  string
	ReportingPlayerId string
	Reason            string
	CreatedAt         time.Time
}
