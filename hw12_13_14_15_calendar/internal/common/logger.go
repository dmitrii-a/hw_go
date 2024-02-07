package common

import (
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/rs/zerolog"
)

// Logger is a main logger for the project(singleton).
var Logger *zerolog.Logger

func setLogger(logLevel string) error {
	level, err := zerolog.ParseLevel(logLevel)
	if IsErr(err) {
		return err
	}
	Logger = logger.InitLogger(level)
	return nil
}

func init() {
	err := setLogger(Config.Server.LogLevel)
	if IsErr(err) {
		panic(err)
	}
}
