package adele

import (
	"log"

	"github.com/cidekar/adele-framework/middleware"
	"github.com/cidekar/adele-framework/mux"
)

type Adele struct {
	AppName    string
	Debug      bool
	ErrorLog   *log.Logger
	InfoLog    *log.Logger
	middleware middleware.Middleware
	Routes     *mux.Mux
	RootPath   string
	Version    string
}
