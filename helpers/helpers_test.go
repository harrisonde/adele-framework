package helpers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/cidekar/adele-framework/render"
	"github.com/cidekar/adele-framework/session"
)

func TestHelpers_Render_Jet(t *testing.T) {

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader("./testdata/views"),
		jet.InDevelopmentMode(),
	)
	var sess = session.Session{
		CookieLifetime: "1",
		CookiePersist:  "true",
		CookieName:     "adele",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}
	testRenderer := &render.Render{
		Directory: "views",
		Renderer:  "jet",
		RootPath:  "./testdata",
		JetViews:  views,
		Session:   sess.InitSession(),
	}

	helpers := &Helpers{Redner: testRenderer}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.Render(w, r, "home", nil, nil)
	})

	middlewareHandler := testRenderer.Session.LoadAndSave(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestHelpers_Render_Go(t *testing.T) {

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader("./testdata/views"),
		jet.InDevelopmentMode(),
	)
	var sess = session.Session{
		CookieLifetime: "1",
		CookiePersist:  "true",
		CookieName:     "adele",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}
	testRenderer := &render.Render{
		Directory: "views",
		Renderer:  "go",
		RootPath:  "./testdata",
		JetViews:  views,
		Session:   sess.InitSession(),
	}

	helpers := &Helpers{Redner: testRenderer}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.Render(w, r, "home", nil, nil)
	})

	middlewareHandler := testRenderer.Session.LoadAndSave(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	fmt.Println(w.Body)
	if w.Code == http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
