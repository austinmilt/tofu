package env

import (
	caarlos "github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type LocalEnv struct {
	c Container
}

// see https://pkg.go.dev/time#ParseDuration (which is what caarl0s uses) for valid duration syntax

type Container struct {
	SammiApiUrl             string `env:"SAMMI_API_URL"`
	SammiAsClientServerPort int    `env:"SAMMI_AS_CLIENT_SERVER_PORT"`
	LogFilePath             string `env:"LOG_FILE_PATH"`
	LogLevel                string `env:"LOG_LEVEL"`
	ProfilerEnabled         bool   `env:"PROFILER_ENABLED"`
	ProfilerPort            int    `env:"PROFILER_PORT"`
}

func (e *LocalEnv) SammiApiUrl() string {
	return e.c.SammiApiUrl
}

func (e *LocalEnv) SammiAsClientServerPort() int {
	return e.c.SammiAsClientServerPort
}

func (e *LocalEnv) LogFilePath() string {
	return e.c.LogFilePath
}

func (e *LocalEnv) LogLevel() string {
	return e.c.LogLevel
}

func (e *LocalEnv) ProfilerEnabled() bool {
	return e.c.ProfilerEnabled
}

func (e *LocalEnv) ProfilerPort() int {
	return e.c.ProfilerPort
}

// Validates that the entire environment is valid, and returns an error
// if not.
func (e *LocalEnv) Validate() error {
	// By this point we've already validated every value.
	return nil
}

func NewLocalEnv() (*LocalEnv, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	var c Container
	err = caarlos.Parse(&c)
	if err != nil {
		return nil, err
	} else {
		return &LocalEnv{c: c}, err
	}
}

func NewLocalEnvWithValues(c Container) LocalEnv {
	return LocalEnv{c: c}
}
