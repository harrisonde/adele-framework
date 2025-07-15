package adele

import (
	"github.com/cidekar/adele-framework/middleware"
	"github.com/cidekar/adele-framework/mux"
	"github.com/sirupsen/logrus"
)

type Adele struct {
	AppName    string
	config     config
	Debug      bool
	Log        *logrus.Logger
	middleware middleware.Middleware
	Routes     *mux.Mux
	RootPath   string
	Version    string
}

type config struct {
	port string
}
