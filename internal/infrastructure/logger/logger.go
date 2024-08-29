package logger

import (
	"github.com/damndelion/test_task_kami/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg configs.Logger) (*zap.SugaredLogger, error) {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(getLogLevel(cfg.LogLevel)),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    getEncoderConfig(),
		ErrorOutputPaths: []string{"stderr"},
		OutputPaths:      []string{"stdout"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func getLogLevel(logLevel string) zapcore.Level {
	// default level
	level := zap.DebugLevel

	switch logLevel {
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.DebugLevel
	}

	return level
}
