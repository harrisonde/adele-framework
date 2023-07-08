package adel

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *Adel) routes() http.Handler {
	mux := chi.NewRouter() // Multiplexer
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if a.Debug {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)
	mux.Use(a.SessionLoad)
	mux.Use(a.NoSurf)

	return mux
}
