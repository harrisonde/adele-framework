package mux

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Cors struct {
	AllowedOrigins []string `yaml:"AllowedOrigins"`

	// AllowOriginFunc is a custom function to validate the origin. It takes the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(r *http.Request, origin string) bool `yaml:"AllowOriginFunc"`

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string `yaml:"AllowedMethods"`

	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string `yaml:"AllowedHeaders"`

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string `yaml:"ExposedHeaders"`

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool `yaml:"AllowCredentials"`

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int `yaml:"MaxAge"`

	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool `yaml:"OptionsPassthrough"`

	// Debugging flag adds additional output to debug server side CORS issues
	Debug bool `yaml:"Debug"`
}

type Mux struct {
	Mux *chi.Mux
}

var Router = &Mux{}

var MuxRouterTree []MuxRouteInfo

type MuxRouteInfo struct {
	Annotation string
	Method     string
	Route      string
	Base       string
	Scope      string
}

type MuxRouteScope struct {
	Scope []string
}
