package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Logger *logrus.Logger
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}
