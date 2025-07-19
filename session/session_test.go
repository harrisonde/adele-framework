package session

import (
	"reflect"
	"testing"

	"github.com/alexedwards/scs/v2"
)

func TestSession_InitSession(t *testing.T) {

	a := &Session{
		CookieLifetime: "100",
		CookiePersist:  "true",
		CookieName:     "Adele",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}

	var sm *scs.SessionManager

	ses := a.InitSession()

	var sessKind reflect.Kind
	var sessType reflect.Type

	rv := reflect.ValueOf(ses)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {

		sessKind = rv.Kind()
		sessType = rv.Type()

		rv = rv.Elem()
	}

	if !rv.IsValid() {
		t.Error("Invalid type or kind:", rv.Kind(), "type:", rv.Type())
	}

	if sessKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind returned while testing cookie. Expected", reflect.ValueOf(sm).Kind(), "and go", sessKind)
	}

	if sessType != reflect.ValueOf(sm).Type() {
		t.Error("wrong type returned while testing cookie. Expected", reflect.ValueOf(sm).Type(), "and go", sessType)
	}
}
