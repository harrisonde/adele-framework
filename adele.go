package adele

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/cidekar/adele-framework/logger"
	"github.com/cidekar/adele-framework/middleware"
	"github.com/cidekar/adele-framework/mux"
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
	a.config = config{
		port:        os.Getenv("PORT"),
		sessionType: os.Getenv("SESSION_TYPE"),
	}

	return nil
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
	fmt.Println("manager", manager)
	return manager, nil
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
