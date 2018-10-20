package sharedlib

import (
	"os"

	"github.com/sirupsen/logrus"
)

// LoggerOption logger option
type LoggerOption struct {
	Level  logrus.Level
	NoLock bool
}

// DefaultLoggerOption default logger option
func DefaultLoggerOption() *LoggerOption {
	return &LoggerOption{
		Level:  logrus.InfoLevel,
		NoLock: true,
	}
}

// NewStdOutLogger init a stdout logger
func NewStdOutLogger(option ...*LoggerOption) *logrus.Logger {
	var loggeroption *LoggerOption
	if len(option) == 0 {
		loggeroption = DefaultLoggerOption()
	} else {
		loggeroption = option[0]
	}

	log := logrus.New()
	log.Out = os.Stdout
	log.SetLevel(loggeroption.Level)
	log.Formatter = new(logrus.JSONFormatter)

	if loggeroption.NoLock {
		log.SetNoLock()
	}

	return log
}
