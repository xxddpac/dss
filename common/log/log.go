package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
)

var (
	lock    sync.Mutex
	_logger *zap.Logger
)

const (
	FormatText           = "text"
	FormatJSON           = "json"
	DefaultLogTimeLayout = "2006-01-02 15:04:05.000"
)

type Config struct {
	LogPath    string
	LogLevel   string
	Compress   bool
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Format     string
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func newLogWriter(logPath string, maxSize, maxBackups, maxAge int, compress bool) io.Writer {
	if logPath == "" || logPath == "-" {
		return os.Stdout
	}
	return &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}

func newZapEncoder() zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(DefaultLogTimeLayout),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	return encoderConfig
}

func newLoggerCore(cfg *Config) zapcore.Core {
	hook := newLogWriter(cfg.LogPath, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress)

	encoderConfig := newZapEncoder()

	atomLevel := zap.NewAtomicLevelAt(getZapLevel(cfg.LogLevel))

	var encoder zapcore.Encoder
	if cfg.Format == FormatJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook)),
		atomLevel,
	)
	return core
}

func newLoggerOptions() []zap.Option {
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	development := zap.Development()
	options := []zap.Option{
		caller,
		callerSkip,
		development,
		zap.Fields(zap.String("project_name", "dss")),
	}
	return options
}

func (c *Config) fillWithDefault() {
	if c.MaxSize <= 0 {
		c.MaxSize = 20
	}
	if c.MaxAge <= 0 {
		c.MaxAge = 7
	}
	if c.MaxBackups <= 0 {
		c.MaxBackups = 7
	}
	if c.LogLevel == "" {
		c.LogLevel = "debug"
	}
	if c.Format == "" {
		c.Format = FormatText
	}
}

func Init(cfg *Config) {
	cfg.fillWithDefault()
	core := newLoggerCore(cfg)
	zapOpts := newLoggerOptions()
	_logger = zap.New(core, zapOpts...)
}

func Logger() *zap.Logger {
	return _logger
}

func Debug(msg string, fields ...zap.Field) {
	_logger.Debug(msg, fields...)
}

func DebugF(msg string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	_logger.Debug(msg)
}

func Info(msg string, fields ...zap.Field) {
	_logger.Info(msg, fields...)
}

func InfoF(msg string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	_logger.Info(msg)
}

func Warn(msg string, fields ...zap.Field) {
	_logger.Warn(msg, fields...)
}

func WarnF(msg string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	_logger.Warn(msg)
}

func Error(msg string, fields ...zap.Field) {
	_logger.Error(msg, fields...)
}

func Errorf(msg string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	_logger.Error(msg)
}

func Panic(msg string, fields ...zap.Field) {
	_logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	_logger.Fatal(msg, fields...)
}
