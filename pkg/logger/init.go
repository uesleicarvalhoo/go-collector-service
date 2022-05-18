package logger

import (
	"os"
	"path"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(cfg Config) error {
	logger, err := newZapLogger(cfg)
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	return nil
}

func newRollingFile(config Config) (zapcore.WriteSyncer, error) {
	if err := os.MkdirAll(config.Directory, os.ModePerm); err != nil {
		Error("Can't create log directory", zap.Error(err), zap.String("path", config.Directory))

		return nil, err
	}

	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		MaxBackups: config.MaxBackups, // files
	}), nil
}

func newConsoleCore(config Config) (zapcore.Core, error) {
	var logLevel zapcore.Level

	if err := logLevel.Set(strings.ToLower(config.ConsoleLevel)); err != nil {
		return nil, err
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
	}

	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})

	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	return zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), levelEnabler), nil
}

func newFileCore(config Config) (zapcore.Core, error) {
	var logLevel zapcore.Level

	if err := logLevel.Set(strings.ToLower(config.FileLevel)); err != nil {
		return nil, err
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
	}

	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel
	})

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	file, err := newRollingFile(config)
	if err != nil {
		return nil, err
	}

	return zapcore.NewCore(encoder, file, levelEnabler), nil
}

func newZapLogger(config Config) (*zap.Logger, error) {
	var cores []zapcore.Core

	if config.ConsoleEnabled {
		consoleCore, err := newConsoleCore(config)
		if err != nil {
			return nil, err
		}

		cores = append(cores, consoleCore)
	}

	if config.FileEnabled {
		fileCore, err := newFileCore(config)
		if err != nil {
			return nil, err
		}

		cores = append(cores, fileCore)
	}

	return zap.New(zapcore.NewTee(cores...)), nil
}
