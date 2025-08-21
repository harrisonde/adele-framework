package render

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cidekar/adele-framework/mux"
)

var pageData = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
	errorMessage  string
}{
	{"go_page", "go", "home", false, "error rendering go template"},
	{"go_page_no_template", "go", "no-file", true, "no error rendering non-existent go template, when on is expected"},
	{"jet_page", "jet", "home", false, "error rendering jet template"},
	{"jet_page_no_template", "jet", "no-file", true, "no error rendering non-existent jet template, when on is expected"},
	{"invalid_renderer_engine", "foo", "home", true, "no error rendering with non-existent template engine"},
}

func TestRender_Page(t *testing.T) {

	r := mux.NewRouter()

	r.Use(testRenderer.Session.LoadAndSave)

	for _, e := range pageData {

		testRenderer.Renderer = e.renderer

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {

			err := testRenderer.Page(w, r, e.template, nil, nil)

			if e.errorExpected {
				if err == nil { // expect error
					t.Errorf("%s: %s:", e.name, e.errorMessage)
				}
			} else {
				if err != nil { // not expecting error
					t.Errorf("%s: %s: %s:", e.name, e.errorMessage, err.Error())
				}
			}
		})

		// If we don't expect an error, execute the request.
		if e.errorExpected != true {
			ts := httptest.NewServer(r)
			defer ts.Close()

			_, body := makeRequest(t, ts, "GET", "/", nil)
			if body != "Adel, let's build something." {
				t.Fatalf("%s", body)
			}

		}

	}
}

func TestRender_GoPage(t *testing.T) {

	testRenderer.Renderer = "go"

	r := mux.NewRouter()

	r.Use(testRenderer.Session.LoadAndSave)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := testRenderer.Page(w, r, "home", nil, nil)
		if err != nil {
			t.Error("Error rendering page", err)
		}
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	if _, body := makeRequest(t, ts, "GET", "/", nil); body != "Adel, let's build something." {
		t.Fatalf("%s", body)
	}
}

func TestRender_JetPage(t *testing.T) {

	testRenderer.Renderer = "jet"

	r := mux.NewRouter()

	r.Use(testRenderer.Session.LoadAndSave)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err := testRenderer.Page(w, r, "home", nil, nil)
		if err != nil {
			t.Error("Error rendering page", err)
		}
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	if _, body := makeRequest(t, ts, "GET", "/", nil); body != "Adel, let's build something." {
		t.Fatalf("%s", body)
	}
}
