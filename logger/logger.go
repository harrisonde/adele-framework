package logger

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type IoToLogWriter struct {
	Entry *logrus.Entry
	Type  string
}

func (w *IoToLogWriter) Write(b []byte) (int, error) {
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}
	if w.Type == "Error" {
		w.Entry.Error(string(b))
	} else {
		w.Entry.Info(string(b))
	}
	return n, nil
}

func NewRequestLogger() func(next http.Handler) http.Handler {
	logger := createLogger()

	return NewStructuredLogger(logger)
}

func createLogger(c ...string) *logrus.Logger {
	l := logrus.New()
	if os.Getenv("LOG_FORMAT") == "JSON" {
		l.SetFormatter(&logrus.JSONFormatter{})
	} else {
		//l.SetFormatter(&logrus.TextFormatter{})
		l.SetFormatter(&Formatter{})
	}
	return l
}
