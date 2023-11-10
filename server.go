package adele

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// Start the web server
func (a *Adele) ListenAndServe() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     a.ErrorLog,
		Handler:      a.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	if a.DB.Pool != nil {
		defer a.DB.Pool.Close()
	}

	if redisPool != nil {
		defer redisPool.Close()
	}

	if badgerPool != nil {
		defer badgerPool.Close()
	}

	a.InfoLog.Printf("\nListening on port %s", os.Getenv("PORT"))

	a.Scheduler.Start()

	return srv.ListenAndServe()
}
