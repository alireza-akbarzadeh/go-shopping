// Package utils provides shared helpers including logging, JWT handling,
// password hashing, consistent API responses, and custom error types.
package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLogger initializes the global logger with the specified level.
// It uses JSON formatting for structured logging suitable for production.
func InitLogger(level string) error {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	Log.SetLevel(lvl)
	return nil
}
