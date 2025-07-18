package middleware

import (
	"fmt"
	"net/http"
)

// Load and save session on each request
func (a *Middleware) SessionLoad(next http.Handler) http.Handler {
	fmt.Println("Middleware SessionLoad a: ", a)
	fmt.Println("Middleware SessionLoad s.Session : ", a.Session)
	return a.Session.LoadAndSave(next)
}
