package adele

import (
	"github.com/alexedwards/scs/v2"
	"github.com/cidekar/adele-framework/mailer"
	"github.com/cidekar/adele-framework/middleware"
	"github.com/cidekar/adele-framework/mux"
	"github.com/sirupsen/logrus"
)

type Adele struct {
	AppName         string
	config          config
	Debug           bool
	Log             *logrus.Logger
	Mail            mailer.Mail
	middleware      middleware.Middleware
	MaintenanceMode bool
	Routes          *mux.Mux
	RootPath        string
	Session         *scs.SessionManager
	Version         string
}

type config struct {
	port        string
	sessionType string
}
