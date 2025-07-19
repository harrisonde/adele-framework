package adele

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Provided a HTTP server implementation that is ready to handle requests.
func (a *Adele) ListenAndServe() error {

	port := "4000"
	if os.Getenv("PORT") != "" {
		os.Getenv("PORT")
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		//TODO: is the error logger necessary here?
		//
		//ErrorLog:     a.ErrorLog,
		Handler:      a.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	a.Log.Debug(fmt.Sprintf("http requests handled on port %s", port))

	return srv.ListenAndServe()
}
