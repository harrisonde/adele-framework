package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cidekar/adele-framework/mux"
	"github.com/cidekar/adele-framework/session"
)

func Test_SessionLoad(t *testing.T) {

	r := mux.NewRouter()

	m := Middleware{}

	var sess = session.Session{
		CookieLifetime: "1",
		CookiePersist:  "true",
		CookieName:     "adele",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}

	m.Session = sess.InitSession()

	r.Use(m.SessionLoad)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/", nil)

	if res.Header.Get("Vary") != "Cookie" {
		t.Error("session middleware was called and did not set the session")
	}

}
