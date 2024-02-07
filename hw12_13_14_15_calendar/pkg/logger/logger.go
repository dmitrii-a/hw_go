package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// InitLogger is a logger initialization.
func InitLogger(level zerolog.Level) *zerolog.Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%s]", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}
	output.FormatErrFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	log := zerolog.New(output).Level(level).With().Timestamp().Caller().Stack().Logger()
	return &log
}
