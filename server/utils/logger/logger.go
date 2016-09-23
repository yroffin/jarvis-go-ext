package logger

import "github.com/Sirupsen/logrus"

// Log : log
var Log *logrus.Logger

// NewLogger : new logger
func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	Log = logrus.New()
	Log.Formatter = new(logrus.TextFormatter)

	return Log
}
