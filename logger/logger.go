package logger

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/hyperits/gosuite/debugger"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	logFile := "logs/app.log"
	rotation := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    32, // 每个日志文件的最大大小（兆字节）
		MaxBackups: 15, // 保留的旧日志文件的最大数量
		MaxAge:     15, // 保留的旧日志文件的最大天数
		Compress:   true,
	}

	log = zerolog.New(rotation).With().Timestamp().Logger()
	SetLevel(InfoLevel)

	// 日志清理
	// 使用 cron 定时清理日志
	c := cron.New()
	// 每天午夜运行清理任务
	_, _ = c.AddFunc("0 0 * * *", func() {
		cleanupLogs(logFile, 15)
	})
	c.Start()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChannel

		// 在完成时关闭日志记录器和 cron 调度器
		_ = rotation.Close()
		c.Stop()
	}()
}

// cleanupLogs 删除早于 `days` 天的日志文件
func cleanupLogs(logFile string, days int) {
	log.Info().Msgf("Start cleanupLogs start with %v within %v days", logFile, days)
	files, err := filepath.Glob(logFile + ".*")
	if err != nil {
		log.Error().Err(err).Msg("Failed to list logfiles")
		return
	}

	for _, file := range files {
		fi, err := os.Stat(file)
		if err != nil {
			log.Error().Err(err).Msg("Failed to stat file info")
			continue
		}

		diff := time.Since(fi.ModTime())
		if diff.Hours() > float64(days*24) {
			err := os.Remove(file)
			if err != nil {
				log.Error().Err(err).Msg("Failed to remove log file")
			} else {
				log.Info().Str("file", file).Msg("Remove log file success")
			}
		}
	}
	log.Info().Msgf("End cleanupLogs start with %v within %v days", logFile, days)
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
