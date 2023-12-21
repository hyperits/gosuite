package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/hyperits/gosuite/debugger"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

// Level defines log levels.
type Level uint8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log = zerolog.New(output).With().Timestamp().Logger()
	SetLevel(DebugLevel)
}

func SetStrLevel(l string) error {
	var err error = nil
	switch l {
	case "debug":
		SetLevel(DebugLevel)
	case "info":
		SetLevel(InfoLevel)
	case "warn":
		SetLevel(WarnLevel)
	case "error":
		SetLevel(ErrorLevel)
	case "fatal":
		SetLevel(FatalLevel)
	case "panic":
		SetLevel(PanicLevel)
	case "no":
		SetLevel(NoLevel)
	case "disable":
		SetLevel(Disabled)
	default:
		err = fmt.Errorf("invalid log level %s", l)
	}
	return err
}

func SetLevel(l Level) {
	zerolog.SetGlobalLevel(zerolog.Level(l))
}

//
// Simple
//

func Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	log.Panic().Msgf(format, v...)
}

//
// Support RuntimeInfo
//

func DebugRTf(info *debugger.RuntimeInfo, format string, v ...interface{}) {
	log.Debug().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func InfoRTf(info *debugger.RuntimeInfo, format string, v ...interface{}) {
	log.Info().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func WarnRTf(info *debugger.RuntimeInfo, format string, v ...interface{}) {
	log.Warn().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func ErrorRTf(info *debugger.RuntimeInfo, format string, v ...interface{}) {
	log.Error().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func FatalRTf(info *debugger.RuntimeInfo, format string, v ...interface{}) {
	log.Fatal().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func PanicRTf(info *debugger.RuntimeInfo, format string, v ...interface{}) {
	log.Panic().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}
