package model

import (
	"sync"
	"time"
)

type Player struct {
	Id           PlayerId
	TwitchId     string
	TwitchHandle string
	AiHandle     string
	GameRecords  map[GameId]PlayerGameRecord
	Takeovers    []Takeover
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Reports      []Report
	Mutex        *sync.Mutex
}

type PlayerId = string

type PlayerGameRecord struct {
	Sessions map[SessionId]PlayerSessionRecord
}

func (p Player) IsNull() bool {
	return (p.Id == "") &&
		(p.TwitchId == "") &&
		(p.TwitchHandle == "") &&
		(p.UpdatedAt == nullTime) &&
		(p.CreatedAt == nullTime) &&
		((p.GameRecords == nil) || (len(p.GameRecords) == 0)) &&
		(len(p.Takeovers) == 0) &&
		(len(p.Reports) == 0)
}

func (p Player) GetLatestOpenSession(gameId GameId) SessionId {
	p.Lock()
	defer p.Unlock()
	var latestSessionId SessionId
	var latestSessionStart time.Time
	for candidateGameId, gameSessions := range p.GameRecords {
		if gameId == candidateGameId {
			for sessionId, session := range gameSessions.Sessions {
				if (session.End == nullTime) && (session.Start.After(latestSessionStart)) {
					latestSessionId = sessionId
					latestSessionStart = session.Start
				}
			}
		}
	}
	return latestSessionId
}

func (p Player) GetLatestTakeover(gameId GameId) *Takeover {
	p.Lock()
	defer p.Unlock()
	var result *Takeover = nil
	for _, candidate := range p.Takeovers {
		if (result == nil) || (candidate.End.After(result.End)) {
			result = &candidate
		}
	}
	return result
}

func (p Player) GetReportsByPlayer() []Report {
	p.Lock()
	defer p.Unlock()
	result := make([]Report, len(p.Reports))
	c := 0
	for _, report := range p.Reports {
		if report.ReportingPlayerId == p.Id {
			result[c] = report
		}
		c += 1
	}
	return result
}

func (p Player) GetReportsAgainstPlayer() []Report {
	p.Lock()
	defer p.Unlock()
	result := make([]Report, len(p.Reports))
	c := 0
	for _, report := range p.Reports {
		if report.ReportedPlayerId == p.Id {
			result[c] = report
		}
		c += 1
	}
	return result
}

func (p Player) GetLatestReportAgainst(other PlayerId) *Report {
	p.Lock()
	defer p.Unlock()
	var result *Report = nil
	for _, report := range p.Reports {
		if (report.ReportingPlayerId == p.Id) &&
			(report.ReportedPlayerId == other) &&
			((result == nil) || report.CreatedAt.After(result.CreatedAt)) {

			result = &report
		}
	}
	return result
}

func (p Player) Lock() {
	p.Mutex.Lock()
}

func (p Player) Unlock() {
	p.Mutex.Unlock()
}

var nullTime = time.Time{}
