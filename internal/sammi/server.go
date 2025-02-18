package sammi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/austinmilt/tofu/internal/env"
	"github.com/austinmilt/tofu/internal/model/visuals"
	"github.com/labstack/echo/v4"
)

func NewServer(
	environment env.Env,
) Server {
	s := &Server{
		ui:             visuals.NewUiState(),
		uiLock:         &sync.Mutex{},
		port:           environment.SammiAsClientServerPort(),
		listeners:      make(map[string]*(func() error)),
		listenersLock:  &sync.Mutex{},
		ioLock:         &sync.Mutex{},
		environment:    environment,
		chatQueue:      make(chan ChatMessageEvent),
		letterCountEnd: time.Unix(0, 0),
		letterCount:    make(map[string]map[string]int),
	}
	return *s
}

type Server struct {
	environment    env.Env
	ui             *visuals.UiState
	uiLock         *sync.Mutex
	port           int
	listeners      map[string]*func() error
	listenersLock  *sync.Mutex
	ioLock         *sync.Mutex
	chatQueue      chan ChatMessageEvent
	letterCountEnd time.Time
	letterCount    map[string]map[string]int
}

func (s *Server) Start() {
	e := echo.New()

	e.POST("/tofu", func(c echo.Context) error {
		request := &Request{}
		if err := c.Bind(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		slog.Info("Got a request", "body", request)
		response := &Response{}
		if request.ChatMessage != nil {
			s.chatQueue <- *request.ChatMessage

		} else if request.StartLetterCount != nil {
			duration, err := time.ParseDuration(request.StartLetterCount.Duration)
			if err != nil {
				c.String(400, "invalid duration")
			}
			s.letterCountEnd = time.Now().Add(duration)
			go func() {
				time.AfterFunc(duration, func() {
					s.sendLetterCountResult()
					clear(s.letterCount)
				})
			}()
		}
		return c.JSON(http.StatusOK, response)
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.handleChatMessages()
	}()

	if err := e.Start(fmt.Sprintf("127.0.0.1:%d", s.port)); (err != nil) && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("error starting web server, panicking", "error", err)
		panic(err)
	}
}

func (s *Server) handleChatMessages() {
	for chatMessage := range s.chatQueue {
		if s.letterCountEnd.After(time.Now()) {
			_, hasChatter := s.letterCount[chatMessage.Chatter]
			if !hasChatter {
				s.letterCount[chatMessage.Chatter] = make(map[string]int)
			}
			for _, r := range chatMessage.Message {
				char := fmt.Sprintf("%c", r)
				_, hasCharacter := s.letterCount[chatMessage.Chatter][char]
				if !hasCharacter {
					s.letterCount[chatMessage.Chatter][char] = 0
				}
				s.letterCount[chatMessage.Chatter][char] = s.letterCount[chatMessage.Chatter][char] + 1
			}
		}
	}
}

func (s *Server) sendLetterCountResult() {
	if len(s.letterCount) == 0 {
		return
	}
	msgBytes, err := json.Marshal(s.letterCount)
	if err != nil {
		slog.Error("error converting letter count to json", "error", err)
	}
	sendChatMessage(string(msgBytes))
}

func sendChatMessage(message string) {
	url := "http://localhost:9450/api"
	body, err := json.Marshal(SammiSetVariableRequest{
		ApiName:       "setVariable",
		VariableName:  "message",
		VariableValue: message,
		ButtonId:      "sendChatListener",
	})
	if err != nil {
		slog.Error("could not marshal message body request", "message", message, "error", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error sending setVariable request", "message", message, "error", err)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	slog.Info("response", "status", resp.Status, "body", response, "error", err)

	if resp.StatusCode >= 400 {
		return
	}

	body, err = json.Marshal(SammiTriggerButtonRequest{
		ApiName:  "triggerButton",
		ButtonId: "sendChatListener",
	})
	if err != nil {
		slog.Error("could not marshal message body request", "message", message, "error", err)
	}
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		slog.Error("error sending triggerButton request", "message", message, "error", err)
	}
	defer resp.Body.Close()

	response, err = io.ReadAll(resp.Body)
	slog.Info("response", "status", resp.Status, "body", response, "error", err)
}

type Response struct {
	Message string `json:"message"`
}

type Request struct {
	StartLetterCount *StartLetterCountRequest `json:"startLetterCount,omitempty"`
	ChatMessage      *ChatMessageEvent        `json:"chatMessage,omitempty"`
}

type StartLetterCountRequest struct {
	Duration string `json:"duration"`
}

type ChatMessageEvent struct {
	Chatter string `json:"chatter"`
	Message string `json:"message"`
}

type SammiSetVariableRequest struct {
	ApiName       string `json:"request"`
	ButtonId      string `json:"buttonID"`
	VariableName  string `json:"name"`
	VariableValue string `json:"value"`
}

type SammiTriggerButtonRequest struct {
	ApiName  string `json:"request"`
	ButtonId string `json:"buttonID"`
}
