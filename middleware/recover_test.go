package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cidekar/adele-framework/mux"
)

func Test_Recover(t *testing.T) {
	// router and middleware
	r := mux.NewRouter()
	m := Middleware{}
	r.Use(m.RecovererWithDebug)

	// output capture for testing
	oldRecovererErrorWriter := recovererErrorWriter
	defer func() {
		recovererErrorWriter = oldRecovererErrorWriter
	}()
	buf := &bytes.Buffer{}
	recovererErrorWriter = buf

	r.Get("/", func(http.ResponseWriter, *http.Request) {
		panic("testing")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "GET", "/", nil)

	t.Helper()
	if !reflect.DeepEqual(res.StatusCode, http.StatusInternalServerError) {
		t.Fatalf("expecting values to be equal but got: '%v' and '%v'", res.StatusCode, http.StatusInternalServerError)
	}
}

func Test_RecoverAbort(t *testing.T) {

	defer func() {
		rcv := recover()

		if rcv != http.ErrAbortHandler {
			t.Fatalf("http.ErrAbortHandler should not be recovered")
		}
	}()

	w := httptest.NewRecorder()

	r := mux.NewRouter()

	m := Middleware{}

	r.Use(m.RecovererWithDebug)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(w, req)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
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
