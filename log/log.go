package log

import (
	"fmt"
	"os"
	"sync"

	"github.com/hyperits/gosuite/kit/debug"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log        zerolog.Logger
	logFile    *lumberjack.Logger
	logFileMux sync.Mutex
)

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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	initLogger()
	SetLevel(InfoLevel)
}

func initLogger() {
	logFile = &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    32, // 每个日志文件的最大大小（兆字节）
		MaxBackups: 15, // 保留的旧日志文件的最大数量
		MaxAge:     15, // 保留的旧日志文件的最大天数
		Compress:   true,
	}

	multi := zerolog.MultiLevelWriter(logFile, os.Stdout)
	log = zerolog.New(multi).With().Timestamp().Logger()
}

// SetLogFileMaxSize 设置日志文件的最大大小（兆字节）。
func SetLogFileMaxSize(sizeMB int) {
	logFileMux.Lock()
	defer logFileMux.Unlock()
	logFile.MaxSize = sizeMB
	initLogger()
}

// SetLogFileMaxBackups 设置保留的旧日志文件的最大数量。
func SetLogFileMaxBackups(backups int) {
	logFileMux.Lock()
	defer logFileMux.Unlock()
	logFile.MaxBackups = backups
	initLogger()
}

// SetLogFileMaxAge 设置保留的旧日志文件的最大天数。
func SetLogFileMaxAge(maxAge int) {
	logFileMux.Lock()
	defer logFileMux.Unlock()
	logFile.MaxAge = maxAge
	initLogger()
}

// SetLogFilePath 设置日志文件的路径。 e.g. logs/app.log
func SetLogFilePath(path string) {
	logFileMux.Lock()
	defer logFileMux.Unlock()
	logFile.Filename = path
	initLogger()
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

func DebugRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	log.Debug().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func InfoRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	log.Info().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func WarnRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	log.Warn().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func ErrorRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	log.Error().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func FatalRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	log.Fatal().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}

func PanicRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	log.Panic().Msgf("[%s:%d %s]: %s", info.File, info.Line, info.Function, fmt.Sprintf(format, v...))
}
