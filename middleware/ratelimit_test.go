package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cidekar/adele-framework/mux"
)

func Test_RateLimiter(t *testing.T) {
	r := mux.NewRouter()
	m := Middleware{}

	r.Use(m.RateLimiter())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/", nil)

	if res.StatusCode != http.StatusOK {
		t.Error("rate limiter middleware returned wrong status code:", res.StatusCode)
	}
}

func Test_RateLimiterCustomValues(t *testing.T) {

	t.Setenv("HTTP_RATE_LIMIT", "1")
	t.Setenv("HTTP_RATE_DURATION", "1")

	r := mux.NewRouter()
	m := Middleware{}

	r.Use(m.RateLimiter())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// make a few requests that will cause the rate limit settings to fail
	testRequest(t, ts, "GET", "/", nil)
	testRequest(t, ts, "GET", "/", nil)
	res, _ := testRequest(t, ts, "GET", "/", nil)

	if res.StatusCode != http.StatusTooManyRequests {
		t.Error("rate limiter middleware returned wrong status code:", res.StatusCode)
	}
}
