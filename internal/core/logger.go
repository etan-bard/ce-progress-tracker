package core

import (
	"go.uber.org/zap"
)

type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	ErrorLevel LogLevel = "error"
)

func MakeLogger(logLevel LogLevel) (*zap.Logger, error) {
	// Initialize Zap logger based on LOG_LEVEL from Viper
	var logger *zap.Logger
	var err error

	switch logLevel {
	case DebugLevel:
		logger, err = zap.NewDevelopment()
		break
	case ErrorLevel:
		logger, err = zap.NewProduction()
		break
	default:
		// Default to debug if LOG_LEVEL is not set or invalid
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}
	return logger, nil
}
