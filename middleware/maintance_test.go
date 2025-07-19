package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cidekar/adele-framework/mux"
)

func Test_CheckForMaintenanceModeMiddleware(t *testing.T) {
	// router and middleware
	r := mux.NewRouter()
	m := Middleware{
		MaintenanceMode: true,
	}

	r.Use(m.CheckForMaintenanceMode)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/", nil)

	if res.StatusCode != http.StatusServiceUnavailable {
		t.Error("check for maintenance mode middleware returned wrong status code:", res.StatusCode)
	}
}

func Test_CheckForMaintenanceModeMiddlewareHealthChecks(t *testing.T) {

	t.Setenv("MAINTENANCE_URL", "/health,/health-check")

	r := mux.NewRouter()

	m := Middleware{
		MaintenanceMode: true,
	}

	r.Use(m.CheckForMaintenanceMode)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/health", nil)

	if res.StatusCode != http.StatusOK {
		t.Error("check for maintenance mode middleware returned wrong status code:", res.StatusCode)
	}

}
