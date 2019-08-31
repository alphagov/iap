package internal

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// SetupLogger will configure the logger in consistent way across the application.
func SetupLogger(json, debug bool) *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	logger.SetOutput(os.Stdout)
	if json {
		logger.SetFormatter(&log.JSONFormatter{})
	}
	if debug {
		logger.SetLevel(log.DebugLevel)
	}

	logger.Debug("debug mode active")

	return logger
}
