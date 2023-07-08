package session

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"reflect"
	"testing"
)

func TestSession_InitSession(t *testing.T) {

	// Setup the session
	a := &Session{
		CookieLifetime: "100",
		CookiePersist:  "true",
		CookieName:     "Adel",
		CookieDomain:   "localhost",
		SessionType:    "cookie",
	}

	// The session manager
	var sm *scs.SessionManager

	// Create the session
	ses := a.InitSession()

	// Checking the return value from the session manager
	var sessKind reflect.Kind
	var sessType reflect.Type

	rv := reflect.ValueOf(ses)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {

		// Print what kinds of things we're getting back
		fmt.Println("For loop:", rv.Kind(), rv.Type(), rv)

		// Assign values from rv to vars
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
