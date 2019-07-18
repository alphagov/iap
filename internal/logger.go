package internal

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetupLogger(debug bool) *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	if debug {
		logger.SetLevel(log.DebugLevel)
	}

	logger.Debug("debug mode active")

	return logger
}
