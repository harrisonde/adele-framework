package adel

import (
	"os"
	"time"

	"net/http"
	"strconv"

	"github.com/go-chi/httprate"
	"github.com/justinas/nosurf"
)

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
