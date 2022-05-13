package logger

import (
	"go.uber.org/zap"
)

// Infof logs a message at level Info on the standard logger.
func Debugf(format string, args ...interface{}) {
	zap.S().Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	zap.S().Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	zap.S().Fatalf(format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	zap.S().Panicf(format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	zap.S().Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	zap.S().Info(args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	zap.S().Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	zap.S().Error(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	zap.S().Fatal(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	zap.S().Panic(args...)
}
