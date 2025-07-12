package logger

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestFormatter_NewStructuredLogger(t *testing.T) {
	l := logrus.Logger{}

	rl := NewStructuredLogger(&l)

	if reflect.TypeOf(rl).String() != "func(http.Handler) http.Handler" {
		t.Error("mux logger did return expected type")
	}
}

func TestFormatter_NewLogEntry(t *testing.T) {
	l := logrus.Logger{}

	sl := StructuredLogger{&l}

	r := http.Request{}

	entry := sl.NewLogEntry(&r)

	if reflect.TypeOf(entry).String() != "*logger.StructuredLoggerEntry" {
		t.Error("mux logger did return expected type of structured logger entry")
	}
}
