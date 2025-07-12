package middleware

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/httprate"
)

type Cookie struct {
	Domain string
	Secure string
}

type Middleware struct {
	Cookie           Cookie
	FrameworkVersion string
	AppName          string
	RootPath         string
	ErrorLog         *log.Logger
	InfoLog          *log.Logger
	MaintenanceMode  bool
	Session          *scs.SessionManager
	Rate             int
	Duration         time.Duration
	Limit            func(requestLimit int, windowLength time.Duration, options ...httprate.Option) func(next http.Handler) http.Handler
}

type FrameworkTrace struct {
	AdeleVersion    string
	AppName         string
	RootPath        string
	FrameCount      int
	GoVersion       string
	FileName        string
	FilePath        string
	PackagePath     string
	MainPath        string
	PanicMessage    string
	PanicType       string
	PanicLine       string
	Stack           []FrameworkTraceEntry
	StackFormatted  []string
	StackRaw        []byte
	SourceRaw       string
	SourceFormatted []string
	SourceHighlight string
}

type FrameworkTraceEntry struct {
	File     string
	Function string
	Line     string
}

// used for testing the recoverer output
var recovererErrorWriter io.Writer = os.Stderr
