package adele

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cidekar/adele-framework/logger"
	"github.com/cidekar/adele-framework/mux"
	"github.com/go-chi/chi/v5/middleware"
	crs "github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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

	log := a.CreateLoggers()

	corsConfig, err := a.LoadCorsConfigurationFromFile(rootPath)
	if err != nil {
		return err
	}

	a.Routes = a.CreateRouter(*corsConfig).(*mux.Mux)
	a.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	a.RootPath = rootPath
	a.Version = Version
	a.Log = log

	return nil
}

// Load the applciation CORS (Cross-Origin Resource Sharing) configuration
// from a YAML file and returning it as a mux.Cors object.
func (a *Adele) LoadCorsConfigurationFromFile(rootPath string) (*mux.Cors, error) {
	configFile, err := os.ReadFile(fmt.Sprintf("%s/config/cors.yml", rootPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read cors config file: %v", err)
	}

	var corsConfig mux.Cors
	err = yaml.Unmarshal(configFile, &corsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &corsConfig, nil
}

// Setup up and configures an HTTP router using the adele mux package. This
// function returns an http.Handler which represents a chain of middleware
// and eventually, the handlers for specific routes.
func (a *Adele) CreateRouter(corsConfig mux.Cors) http.Handler {
	mux := mux.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
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

	if a.Debug {
		muxLog := logrus.New()
		if os.Getenv("LOG_FORMAT") == "JSON" {
			muxLog.SetFormatter(&logrus.JSONFormatter{})
		} else {
			muxLog.SetFormatter(&logrus.TextFormatter{})
		}
		mux.Use(logger.NewStructuredLogger(muxLog))
		mux.Use(a.middleware.RecovererWithDebug)
	} else {
		mux.Use(middleware.Recoverer)
	}

	mux.Use(a.middleware.SessionLoad)
	mux.Use(a.middleware.CheckForMaintenanceMode)

	return mux
}

// Create application loggers
// func (a *Adele) CreateLoggers() (*log.Logger, *log.Logger) {
func (a *Adele) CreateLoggers() *logrus.Logger {

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	if os.Getenv("LOG_FORMAT") == "JSON" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	return log
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
