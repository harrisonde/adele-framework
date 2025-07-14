package logger

import (
	"github.com/sirupsen/logrus"
)

type IoToLogWriter struct {
	Entry *logrus.Entry
	Type  string
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}
