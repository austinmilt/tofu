package env

type Env interface {
	// Validates that the entire environment is valid, and returns an error
	// if not.
	Validate() error

	// URL of the SAMMI server API server that SAMMI is listening to
	// like http://localhost:9450/api
	SammiApiUrl() string

	// Port that the local SAMMI program should use to make requests to tofu (this app)
	SammiAsClientServerPort() int

	// path relative to the project directory to create the log output file
	LogFilePath() string

	// TODO change to enum
	// log level (DEBUG < INFO < WARN < ERROR)
	LogLevel() string

	// whether to run pprof to help profile the app
	ProfilerEnabled() bool

	// the port that the pprofiler server should listen to
	ProfilerPort() int
}
