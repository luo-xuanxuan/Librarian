package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger defines the interface for the logging operations. This can be extended with more methods as needed.
type Logger interface {
	WithField(key string, value interface{}) Logger
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

// logWrapper implements the Logger interface using Logrus
type logWrapper struct {
	*logrus.Entry
}

var Log *logWrapper

// New creates and returns a Logger configured with default settings.
func init() {
	var baseLogger = logrus.New()
	baseLogger.Out = os.Stdout
	baseLogger.Formatter = &logrus.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	}
	baseLogger.Level = logrus.DebugLevel // Set the default level. You can make this configurable.

	Log = &logWrapper{baseLogger.WithFields(logrus.Fields{})}
}

func (l *logWrapper) WithField(key string, value interface{}) Logger {
	return &logWrapper{l.Entry.WithField(key, value)}
}

// Debug logs a message at level Debug on the standard logger.
func (l *logWrapper) Debug(args ...interface{}) {
	l.Entry.Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func (l *logWrapper) Info(args ...interface{}) {
	l.Entry.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l *logWrapper) Warn(args ...interface{}) {
	l.Entry.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func (l *logWrapper) Error(args ...interface{}) {
	l.Entry.Error(args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l *logWrapper) Fatal(args ...interface{}) {
	l.Entry.Fatal(args...)
}
