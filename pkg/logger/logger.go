package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

// Initialize sets up the logger with the specified configuration
func Initialize(logLevel string, prettyPrint bool) {
	// Set the global time format
	zerolog.TimeFieldFormat = time.RFC3339

	// Set the global log level
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Configure the logger output
	if prettyPrint {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Replace the global logger
	log.Logger = logger
}

// GetLogger returns the configured logger instance
func GetLogger() zerolog.Logger {
	// If logger hasn't been initialized, use a default configuration
	if logger.GetLevel() == zerolog.NoLevel {
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info"
		}
		prettyPrint := os.Getenv("LOG_PRETTY") == "true"
		Initialize(logLevel, prettyPrint)
	}
	return logger
}
