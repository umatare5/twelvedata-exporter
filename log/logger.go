// Package log provides a simple logging interface.
package log

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// Info logs a message at level Info.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof logs a message at level Info.
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Errorf logs a message at level Error.
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal logs a message at level Fatal.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
