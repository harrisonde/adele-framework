package adel

import (
	"net/http"

	"github.com/harrisonde/adel/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Adel) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if a.Debug {
		mux.Use(logger.NewRequestLogger())
	}

	mux.Use(middleware.Recoverer)
	mux.Use(a.SessionLoad)
	mux.Use(a.NoSurf)
	mux.Use(a.CheckForMaintenanceMode)

	return mux
}
