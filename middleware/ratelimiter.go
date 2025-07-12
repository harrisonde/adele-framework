package middleware

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/httprate"
)

func (a *Middleware) RateLimiter() func(next http.Handler) http.Handler {
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
