package logger

import (
	"github.com/sirupsen/logrus"
)

type IoToLogWriter struct {
	Entry *logrus.Entry
	Type  string
}
