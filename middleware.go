package adel

import (
	"fmt"
	"os"

	"net/http"
	"strconv"
	"strings"

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

	// Pattern of string that does not get a token
	// https://github.com/justinas/nosurf/blob/master/exempt.go#L55
	// ExemptPath(), ExemptGlob(), ExemptGlobs() ...
	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   a.config.cookie.domain,
	})

	return csrfHandler
}

func (a *Adel) CheckForMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maintenanceMode {
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
