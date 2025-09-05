package render

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CloudyKit/jet/v6"
	"github.com/cidekar/adele-framework/session"
)

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

var testRenderer = Render{
	Directory: "views",
	Renderer:  "",
	RootPath:  "./testdata",
	JetViews:  views,
	Session:   sess.InitSession(),
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func makeRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
