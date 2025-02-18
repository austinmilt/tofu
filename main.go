package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/austinmilt/tofu/internal/env"
	"github.com/austinmilt/tofu/internal/sammi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// how to get the user access token:
// manual token, go here https://twitchtokengenerator.com
// or use this parameterized url
// https://twitchtokengenerator.com/quick/sNZSEZzznJ

// https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/#authorization-code-grant-flow
// scopes at https://dev.twitch.tv/docs/authentication/scopes/#twitch-access-token-scopes
// or using the twitch cli https://dev.twitch.tv/docs/cli/token-command/

func main() {
	var wg sync.WaitGroup

	environment := initEnv()
	initLogging(environment)

	// sammiClient := sammi.NewClient(environment)
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	err := sammiClient.Start()
	// 	if err != nil {
	// 		slog.Error("error starting SAMMI websocket client", "error", err)
	// 		panic(err)
	// 	}
	// }()

	sammiServer := sammi.NewServer(environment)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sammiServer.Start()
	}()

	// program profiler
	// https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/
	if environment.ProfilerEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			http.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
			http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", environment.ProfilerPort()), nil)
		}()
	}

	wg.Wait()
}

func initEnv() env.Env {
	environment, err := env.NewLocalEnv()
	if err != nil {
		panic(err)
	}
	return environment
}

func initLogging(environment env.Env) {
	log.SetOutput(&lumberjack.Logger{
		Filename:   environment.LogFilePath(),
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   false,
	})

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var logLevel slog.Level
	switch environment.LogLevel() {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		panic("invalid log level " + environment.LogLevel())
	}
	h := slog.NewTextHandler(log.Writer(), &slog.HandlerOptions{
		Level: logLevel,
	})
	slog.SetDefault(slog.New(h))
}

type startFunc = func()
type stopFunc = func()
