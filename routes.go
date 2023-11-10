package adele

import (
	"net/http"

	"github.com/harrisonde/adele/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Adele) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(a.rateLimiter())

	if a.Debug {
		mux.Use(logger.NewRequestLogger())
	}

	mux.Use(middleware.Recoverer)
	mux.Use(a.SessionLoad)
	mux.Use(a.CheckForMaintenanceMode)

	return mux
}
