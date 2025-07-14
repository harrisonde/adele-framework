package adele

import (
	"github.com/cidekar/adele-framework/middleware"
	"github.com/cidekar/adele-framework/mux"
	"github.com/sirupsen/logrus"
)

type Adele struct {
	AppName string
	Debug   bool
	// ErrorLog   *log.Logger
	// InfoLog    *log.Logger
	Log        *logrus.Logger
	middleware middleware.Middleware
	Routes     *mux.Mux
	RootPath   string
	Version    string
}
