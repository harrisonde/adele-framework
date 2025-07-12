package adele

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Provided a HTTP server implementation that is ready to handle requests.
func (a *Adele) ListenAndServe() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     a.ErrorLog,
		Handler:      nil,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	a.InfoLog.Printf("adele is ready to handle http requests")

	return srv.ListenAndServe()
}
