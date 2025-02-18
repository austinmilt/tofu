package visuals

func NewUiState() *UiState {
	defThreshold := float64(1)
	defTopPlayers := []TopPlayer{}
	defActivePlayers := 0
	return &UiState{
		InputWindow: InputWindow{
			Inputs:             map[string]float64{},
			ConsensusThreshold: defThreshold,
		},
		ActiveTakeover:    nil,
		TopPlayers:        defTopPlayers,
		ActivePlayerCount: defActivePlayers,
		ChatLog:           []ChatMessage{},
		EventLog:          []Event{},
	}
}

type UiState struct {
	InputWindow       InputWindow              `json:"inputWindow"`
	ActiveTakeover    *Takeover                `json:"takeover,omitempty"`
	TopPlayers        []TopPlayer              `json:"topPlayers"`
	ActivePlayerCount int                      `json:"activePlayerCount"`
	ChatLog           []ChatMessage            `json:"chatLog"`
	EventLog          []Event                  `json:"eventLog"`
	PcPosition        *PlayerCharacterPosition `json:"pcPosition"`
}

type InputWindow struct {
	Inputs             map[string]float64 `json:"inputs"`
	ConsensusThreshold float64            `json:"consensusThreshold"`
}

type TopPlayer struct {
	Handle         string `json:"handle"`
	ConsensusVotes int    `json:"consensusVotes"`
	Votes          int    `json:"votes"`
}

type Takeover struct {
	PlayerHandle     string `json:"playerHandle"`
	StartTimeRfc3339 string `json:"startTimeRfc3339"`
	EndTimeRfc3339   string `json:"endTimeRfc3339"`
}

type ChatMessage struct {
	ChatterHandle string `json:"handle"`
	Message       string `json:"message"`
}

type Event = string

type PlayerCharacterPosition struct {
	Row int `json:"row"`
	Col int `json:"col"`
}
