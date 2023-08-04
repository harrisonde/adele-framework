package logger

import (
	"bytes"
	"fmt"
	"runtime"

	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type colorPallet struct {
	Background int
	Foreground int
}

type Formatter struct {
	CallerFirst           bool
	CustomCallerFormatter func(*runtime.Frame) string
	FieldsOrder           []string
	HideKeys              bool
	NoColors              bool
	NoUppercaseLevel      bool
	NoFieldsColors        bool
	ShowFullLevel         bool
	TimestampFormat       string
	TrimMessages          bool
	Component             string
}

const escape = "\x1b"

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

	ts := f.TimestampFormat
	if ts == "" {
		ts = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	// write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if f.CallerFirst {
		f.writeCaller(b, entry)
	}

	// Pad the output for better readability
	fmt.Printf(" ")

	if !f.NoColors {
		levelColor := getColorByLevel(entry.Level)
		fmt.Printf("%s[%dm%s[%dm", escape, levelColor.Foreground, escape, levelColor.Background)
	}

	// Pad the output for better readability
	b.WriteString(" ")

	// Start log level tag
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}

	b.WriteString(" ")

	if !f.NoColors {
		b.WriteString("\x1b[0m")
	}

	if f.Component != "" {
		b.WriteString(" ")
		b.WriteString(f.Component + ":")

	}
	// Pad the output for better readability
	b.WriteString(" ")

	// write time
	b.WriteString(entry.Time.Format(ts))

	b.WriteString(" ")

	// write fields
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	// write message
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	if !f.CallerFirst {
		f.writeCaller(b, entry)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				" (%s:%d %s)",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)

	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	b.WriteString(" ")

	if f.HideKeys {
		fmt.Fprintf(b, "[%v]", entry.Data[field])
	} else {
		fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
	}

}

func getColorByLevel(level logrus.Level) colorPallet {
	var colorRed = colorPallet{
		Background: 41,
		Foreground: 37,
	}

	var colorYellow = colorPallet{
		Background: 43,
		Foreground: 37,
	}

	var colorBlue = colorPallet{
		Background: 44,
		Foreground: 37,
	}

	var colorWhite = colorPallet{
		Background: 0,
		Foreground: 37,
	}

	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorWhite
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
