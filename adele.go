package adele

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/cidekar/adele-framework/logger"
	"github.com/cidekar/adele-framework/mailer"
	"github.com/cidekar/adele-framework/middleware"
	"github.com/cidekar/adele-framework/mux"
	"github.com/cidekar/adele-framework/render"
	"github.com/cidekar/adele-framework/session"
	crs "github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

const Version = "v0.0.0"

// Create a new instance of the Adele type using a pointer to Adele with the
// root path of the application as a argument. The new-up is called by project adele's consuming package
// to bootstrap the framework.
func (a *Adele) New(rootPath string) error {

	directories := []string{"data", "handlers", "logs", "jobs", "middleware", "migrations", "public", "resources", "resources/views", "resources/mail", "tmp", "screenshots"}

	err := a.CreateDirectories(rootPath, directories)
	if err != nil {
		return err
	}

	err = a.CreateEnvironmentFile(rootPath)
	if err != nil {
		return err
	}

	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	a.Log = logger.CreateLogger()

	sess, err := a.BootstrapSessionManager()
	if err != nil {
		return err
	}

	a.Session = sess

	a.BootstrapMiddleware()

	muxRouter, err := a.BootstrapMux(rootPath)
	if err != nil {
		return err
	}

	a.Routes = muxRouter.(*mux.Mux)
	a.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	a.RootPath = rootPath
	a.Version = Version
	a.ViewsTemplateDir = Getenv("VIEWS_TEMPLATE_DIR", "resources/views")
	a.config = config{
		port:        os.Getenv("PORT"),
		renderer:    Getenv("RENDERER", "jet"),
		sessionType: os.Getenv("SESSION_TYPE"),
	}

	a.Mail = a.BoootstrapMailer()

	a.JetViews = a.BootstrapJetEngine()

	a.Render = a.BootstrapRender()

	return nil
}

// Configure the mailer for the application by initializing mailer struct. The mailer
// values are populated by the environemnt variables parsed from the .env file
// at the root of the application.
func (a *Adele) BoootstrapMailer() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	m := mailer.Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Templates:   a.RootPath + "/resources/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		Jobs:        make(chan mailer.Message, 20),
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		APIKey:      os.Getenv("MAILER_KEY"),
		APIUrl:      os.Getenv("MAILER_URL"),
	}
	return m
}

// Configure the middleware for the application by initializing a middleware struct,
// populating its values using the application configuration.
func (a *Adele) BootstrapMiddleware() {
	myMiddleware := middleware.Middleware{
		FrameworkVersion: a.Version,
		AppName:          a.AppName,
		RootPath:         a.RootPath,
		Log:              a.Log,
		Session:          a.Session,
		MaintenanceMode:  a.MaintenanceMode,
	}

	a.middleware = myMiddleware
}

// Configure and create the session manager by initializing a session struct, populating
// its cookie fields by retrieving values from environment variables.
func (a *Adele) BootstrapSessionManager() (*scs.SessionManager, error) {

	session := session.Session{
		CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
		CookieLifetime: os.Getenv("COOKIE_LIFETIME"),
		CookieName:     os.Getenv("COOKIE_NAME"),
		CookiePersist:  os.Getenv("COOKIE_PERSIST"),
		CookieSecure:   os.Getenv("COOKIE_SECURE"),
	}

	switch strings.ToLower(os.Getenv("SESSION_TYPE")) {
	case "redis":
		//...

	case "mysql", "postgres", "mariadb", "postgresql":
		//...
	default:
		a.Log.Warn("sessions using in-memory session store")
	}

	manager := session.InitSession()
	return manager, nil
}

// Setup This code is setting up the Jet template engine for your Adele framework with
// different configurations based on whether the application is in debug/development mode
// or production mode—enables features that help during development but would hurt performance
// in production (like not caching templates and reloading them on every request).
func (a *Adele) BootstrapJetEngine() *jet.Set {
	loader := jet.NewOSFileSystemLoader(fmt.Sprintf("%s/%s", a.RootPath, a.ViewsTemplateDir))

	var views *jet.Set
	if a.Debug {
		views = jet.NewSet(loader, jet.InDevelopmentMode())
	} else {
		views = jet.NewSet(loader)
	}

	views.AddGlobal("APP_DEBUG", a.Debug)

	return views
}

// Setup and configuring a render engine- initializes a rendering system that handles
// template rendering for web responses (HTML pages, emails, etc.). The render system
// handles Rendering HTML templates for web pages, passing session data to templates,
// and, managing template inheritance and layouts.
func (a *Adele) BootstrapRender() *render.Render {
	r := render.Render{
		Directory: a.ViewsTemplateDir,
		Renderer:  a.config.renderer,
		RootPath:  a.RootPath,
		Port:      a.config.port,
		JetViews:  a.JetViews,
		Session:   a.Session,
	}

	return &r
}

// Setup up and configures an HTTP router using the adele mux package. This
// function returns an http.Handler which represents a chain of middleware
// and eventually, the handlers for specific routes.
func (a *Adele) BootstrapMux(rootPath string) (http.Handler, error) {

	// Load the applciation CORS (Cross-Origin Resource Sharing) configuration
	// from a YAML file and returning it as a mux.Cors object.
	configFile, err := os.ReadFile(fmt.Sprintf("%s/config/cors.yml", rootPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read cors config file: %v", err)
	}

	var corsConfig mux.Cors
	err = yaml.Unmarshal(configFile, &corsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	mux := mux.NewRouter()
	mux.Use(middleware.RequestID())
	mux.Use(middleware.RealIP())
	mux.Use(a.middleware.RateLimiter())

	corsOptions := crs.Options{
		AllowedOrigins:   corsConfig.AllowedOrigins,
		AllowedMethods:   corsConfig.AllowedMethods,
		AllowedHeaders:   corsConfig.AllowedHeaders,
		ExposedHeaders:   corsConfig.ExposedHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           corsConfig.MaxAge,
	}
	mux.Use(crs.Handler(corsOptions))

	debugMode, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	if debugMode {
		mux.Use(logger.HttpRequesLogger(logger.CreateLogger()))
		mux.Use(a.middleware.RecovererWithDebug)
	} else {
		mux.Use(middleware.Recoverer())
	}

	mux.Use(a.middleware.SessionLoad)
	mux.Use(a.middleware.CheckForMaintenanceMode)

	return mux, nil
}

// Ensure that a environment file at a specific path exists, creating it if it's missing, and returning
// any errors that may arise.
func (a *Adele) CreateEnvironmentFile(rootPath string) error {
	err := a.CreateFileIfNotExist(fmt.Sprintf("%s/.env", rootPath))
	if err != nil {
		return err
	}
	return nil
}

// Create all nonexistent parent directories
func (a *Adele) CreateDirectories(rootPath string, directories []string) error {
	for _, path := range directories {
		err := a.CreateDirIfNotExist(rootPath + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}
