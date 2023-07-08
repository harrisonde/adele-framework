package adel

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

// Load and save session on each request
func (a *Adel) SessionLoad(next http.Handler) http.Handler {
	a.InfoLog.Println("SessionLoad called")
	return a.Session.LoadAndSave(next)
}

// Setup and return CSRF token setup
func (a *Adel) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(a.config.cookie.secure)

	// TO DO
	//
	// How will this fit into API request since they do not need csrf?
	//
	// Be nice to do this on the route level and not the handler.
	//
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
