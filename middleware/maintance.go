package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func (a *Middleware) CheckForMaintenanceMode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if a.MaintenanceMode {
			// urls accessible while application is in maintenance mode e.g., health check url.
			urls := strings.Split(os.Getenv("MAINTENANCE_URL"), ",")
			if len(urls) > 1 {
				for _, url := range urls {
					if strings.Contains(r.URL.Path, url) {
						return
					}
				}
			}

			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Retry-After:", "300")
			w.Header().Set("Cache-Control:", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
			http.ServeFile(w, r, fmt.Sprintf("%s/public/maintenance.html", a.RootPath))
			return
		}
		next.ServeHTTP(w, r)
	})
}
