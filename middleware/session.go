package middleware

import (
	"net/http"
)

// Load and save session on each request
func (a *Middleware) SessionLoad(next http.Handler) http.Handler {
	return a.Session.LoadAndSave(next)
}
