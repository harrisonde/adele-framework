package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Create page data that is a slice of struct
// And use a table test
var pageData = []struct {
	name          string // Name of test
	renderer      string // go or jet
	template      string // the name of the template
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

	for _, e := range pageData {
		r, err := http.NewRequest("GET", "/url", nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()

		testRenderer.Renderer = e.renderer
		testRenderer.RootPath = "./testdata"

		err = testRenderer.Page(w, r, e.template, nil, nil)
		if e.errorExpected {
			if err == nil { // expect error
				t.Errorf("%s: %s:", e.name, e.errorMessage)
			}
		} else {
			if err != nil { // not expecting error
				t.Errorf("%s: %s: %s:", e.name, e.errorMessage, err.Error())
			}
		}

	}
}

func TestRender_GoPage(t *testing.T) {
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"

	// Go page
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering page", err)
	}

}

func TestRender_JetPage(t *testing.T) {
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	testRenderer.Renderer = "jet"

	// Go page
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering page", err)
	}
}
