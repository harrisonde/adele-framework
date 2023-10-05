package adel

import (
	"fmt"
	"os"
	"strings"
	"time"

	"net/http"
	"strconv"

	"github.com/go-chi/httprate"
	"github.com/justinas/nosurf"
)

func (a *Adel) CheckForMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.MaintenanceMode {
			isAccessible := false
			urls := strings.Split(os.Getenv("MAINTENANCE_URL"), ",")
			for _, url := range urls {
				if strings.Contains(r.URL.Path, url) {
					isAccessible = true
					return
				}
			}
			if !isAccessible {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Header().Set("Retry-After:", "300")
				w.Header().Set("Cache-Control:", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
				http.ServeFile(w, r, fmt.Sprintf("%s/public/maintenance.html", a.RootPath))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Load and save session on each request
func (a *Adel) SessionLoad(next http.Handler) http.Handler {
	return a.Session.LoadAndSave(next)
}

// Setup and return CSRF token setup
func (a *Adel) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(a.config.cookie.secure)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   a.config.cookie.domain,
	})

	return csrfHandler
}

func (a *Adel) rateLimiter() func(next http.Handler) http.Handler {
	var rate int
	var duration int

	// Default to 100 requests per minute
	rateDefault := 100
	durationDefault := 1

	_, ok := os.LookupEnv("HTTP_RATE_LIMIT")
	if ok {
		r, err := strconv.Atoi(os.Getenv("HTTP_RATE_LIMIT"))
		if err != nil {
			rate = rateDefault
		} else {
			rate = r
		}
	} else {
		rate = rateDefault
	}

	_, ok = os.LookupEnv("HTTP_RATE_DURATION")
	if ok {
		r, err := strconv.Atoi(os.Getenv("HTTP_RATE_DURATION"))
		if err != nil {
			duration = durationDefault
		} else {
			duration = r
		}
	} else {
		duration = durationDefault
	}

	return httprate.LimitByIP(rate, time.Duration(duration)*time.Minute)
}
