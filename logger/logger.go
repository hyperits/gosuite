package logger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hyperits/gosuite/kit/debug"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger     zerolog.Logger
	logFile    *lumberjack.Logger
	logFileMux sync.Mutex
	configured bool
)

// Level 日志级别
type Level uint8

const (
	// DebugLevel 调试级别
	DebugLevel Level = iota
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 致命错误级别
	FatalLevel
	// PanicLevel panic 级别
	PanicLevel
	// NoLevel 无级别
	NoLevel
	// Disabled 禁用日志
	Disabled
)

// Config 日志配置
type Config struct {
	// 日志文件路径，为空则只输出到控制台
	FilePath string
	// 日志文件最大大小（MB），默认 32
	MaxSize int
	// 保留的旧日志文件最大数量，默认 15
	MaxBackups int
	// 保留的旧日志文件最大天数，默认 15
	MaxAge int
	// 是否压缩旧日志文件，默认 true
	Compress bool
	// 日志级别，默认 InfoLevel
	Level Level
	// 是否输出到控制台，默认 true
	Console bool
	// 是否输出调用者信息，默认 false
	Caller bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		FilePath:   "",
		MaxSize:    32,
		MaxBackups: 15,
		MaxAge:     15,
		Compress:   true,
		Level:      InfoLevel,
		Console:    true,
		Caller:     false,
	}
}

func init() {
	// 默认初始化：仅控制台输出
	initDefault()
}

// initDefault 默认初始化（仅控制台）
func initDefault() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Init 使用配置初始化日志
func Init(cfg *Config) {
	logFileMux.Lock()
	defer logFileMux.Unlock()

	if cfg == nil {
		cfg = DefaultConfig()
	}

	var writers []io.Writer

	// 配置文件输出
	if cfg.FilePath != "" {
		logFile = &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		writers = append(writers, logFile)
	}

	// 配置控制台输出
	if cfg.Console {
		writers = append(writers, os.Stdout)
	}

	// 如果没有任何输出，默认输出到控制台
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	// 创建多输出 writer
	multi := zerolog.MultiLevelWriter(writers...)

	// 创建 logger
	ctx := zerolog.New(multi).With().Timestamp()
	if cfg.Caller {
		ctx = ctx.Caller()
	}
	logger = ctx.Logger()

	// 设置日志级别
	zerolog.SetGlobalLevel(zerolog.Level(cfg.Level))
	configured = true
}

// SetLogFileMaxSize 设置日志文件的最大大小（兆字节）
// Deprecated: 请使用 Init 方法配置
func SetLogFileMaxSize(sizeMB int) {
	logFileMux.Lock()
	defer logFileMux.Unlock()
	if logFile != nil {
		logFile.MaxSize = sizeMB
	}
}

// SetLogFileMaxBackups 设置保留的旧日志文件的最大数量
// Deprecated: 请使用 Init 方法配置
func SetLogFileMaxBackups(backups int) {
	logFileMux.Lock()
	defer logFileMux.Unlock()
	if logFile != nil {
		logFile.MaxBackups = backups
	}
}

// SetLogFileMaxAge 设置保留的旧日志文件的最大天数
// Deprecated: 请使用 Init 方法配置
func SetLogFileMaxAge(maxAge int) {
	logFileMux.Lock()
	defer logFileMux.Unlock()
	if logFile != nil {
		logFile.MaxAge = maxAge
	}
}

// SetLogFilePath 设置日志文件的路径
// Deprecated: 请使用 Init 方法配置
func SetLogFilePath(path string) {
	Init(&Config{
		FilePath:   path,
		MaxSize:    32,
		MaxBackups: 15,
		MaxAge:     15,
		Compress:   true,
		Level:      InfoLevel,
		Console:    true,
	})
}

// SetStrLevel 通过字符串设置日志级别
func SetStrLevel(l string) error {
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
	case "disable", "disabled":
		SetLevel(Disabled)
	default:
		return fmt.Errorf("invalid log level: %s", l)
	}
	return nil
}

// SetLevel 设置日志级别
func SetLevel(l Level) {
	zerolog.SetGlobalLevel(zerolog.Level(l))
}

//
// 基础日志方法
//

// Debugf 输出调试日志
func Debugf(format string, v ...interface{}) {
	logger.Debug().Msgf(format, v...)
}

// Infof 输出信息日志
func Infof(format string, v ...interface{}) {
	logger.Info().Msgf(format, v...)
}

// Warnf 输出警告日志
func Warnf(format string, v ...interface{}) {
	logger.Warn().Msgf(format, v...)
}

// Errorf 输出错误日志
func Errorf(format string, v ...interface{}) {
	logger.Error().Msgf(format, v...)
}

// Fatalf 输出致命错误日志并退出程序
func Fatalf(format string, v ...interface{}) {
	logger.Fatal().Msgf(format, v...)
}

// Panicf 输出 panic 日志并触发 panic
func Panicf(format string, v ...interface{}) {
	logger.Panic().Msgf(format, v...)
}

//
// 带运行时信息的日志方法
//

// DebugRTf 输出带运行时信息的调试日志
func DebugRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	logger.Debug().
		Str("file", info.File).
		Int("line", info.Line).
		Str("func", info.Function).
		Msgf(format, v...)
}

// InfoRTf 输出带运行时信息的信息日志
func InfoRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	logger.Info().
		Str("file", info.File).
		Int("line", info.Line).
		Str("func", info.Function).
		Msgf(format, v...)
}

// WarnRTf 输出带运行时信息的警告日志
func WarnRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	logger.Warn().
		Str("file", info.File).
		Int("line", info.Line).
		Str("func", info.Function).
		Msgf(format, v...)
}

// ErrorRTf 输出带运行时信息的错误日志
func ErrorRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	logger.Error().
		Str("file", info.File).
		Int("line", info.Line).
		Str("func", info.Function).
		Msgf(format, v...)
}

// FatalRTf 输出带运行时信息的致命错误日志并退出程序
func FatalRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	logger.Fatal().
		Str("file", info.File).
		Int("line", info.Line).
		Str("func", info.Function).
		Msgf(format, v...)
}

// PanicRTf 输出带运行时信息的 panic 日志并触发 panic
func PanicRTf(info *debug.RuntimeInfo, format string, v ...interface{}) {
	logger.Panic().
		Str("file", info.File).
		Int("line", info.Line).
		Str("func", info.Function).
		Msgf(format, v...)
}

//
// 结构化日志方法
//

// WithFields 返回带有字段的日志事件
func WithFields(fields map[string]interface{}) *zerolog.Event {
	event := logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	return event
}

// Logger 返回底层的 zerolog.Logger
func Logger() zerolog.Logger {
	return logger
}

