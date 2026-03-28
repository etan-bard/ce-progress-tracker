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
	switch logLevel {
	case DebugLevel:
		return zap.NewDevelopment()
	case ErrorLevel:
		return zap.NewProduction()
	default:
		return zap.NewDevelopment()
	}
}
