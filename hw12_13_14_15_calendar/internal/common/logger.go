package common

import (
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/pkg/logger"
	"github.com/rs/zerolog"
)

// Logger is a main logger for the project(singleton).
var Logger *zerolog.Logger

func init() {
	level, err := zerolog.ParseLevel(Config.Server.LogLevel)
	if IsErr(err) {
		panic(err)
	}
	Logger = logger.InitLogger(level)
}
