package adele

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/cidekar/adele-framework/helpers"
	"github.com/cidekar/adele-framework/mailer"
	"github.com/cidekar/adele-framework/middleware"
	"github.com/cidekar/adele-framework/mux"
	"github.com/cidekar/adele-framework/render"
	"github.com/sirupsen/logrus"
)

type Adele struct {
	AppName          string
	config           config
	Debug            bool
	Helpers          *helpers.Helpers
	JetViews         *jet.Set
	Log              *logrus.Logger
	Mail             mailer.Mail
	middleware       middleware.Middleware
	MaintenanceMode  bool
	Render           *render.Render
	Routes           *mux.Mux
	RootPath         string
	Session          *scs.SessionManager
	Version          string
	ViewsTemplateDir string
}

type config struct {
	port        string
	renderer    string
	sessionType string
}
