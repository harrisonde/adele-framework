package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
)

func TestRealIP(t *testing.T) {
	// Create middleware
	middleware := RealIP()
	if middleware == nil {
		t.Fatal("RealIP() returned nil middleware")
	}

	// Test handler that captures RemoteAddr
	var capturedRemoteAddr string
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRemoteAddr = r.RemoteAddr
	})

	// Wrap test handler with middleware
	wrappedHandler := middleware(testHandler)

	tests := []struct {
		name           string
		headers        map[string]string
		originalRemote string
		expectedRemote string
	}{
		{
			name: "True-Client-IP header (highest priority)",
			headers: map[string]string{
				"True-Client-IP":  "192.168.1.100",
				"X-Real-IP":       "10.0.0.1",
				"X-Forwarded-For": "172.16.0.1",
			},
			originalRemote: "127.0.0.1:8080",
			expectedRemote: "192.168.1.100",
		},
		{
			name: "X-Real-IP header (second priority)",
			headers: map[string]string{
				"X-Real-IP":       "10.0.0.1",
				"X-Forwarded-For": "172.16.0.1",
			},
			originalRemote: "127.0.0.1:8080",
			expectedRemote: "10.0.0.1",
		},
		{
			name: "X-Forwarded-For header (third priority)",
			headers: map[string]string{
				"X-Forwarded-For": "172.16.0.1, 192.168.1.1",
			},
			originalRemote: "127.0.0.1:8080",
			expectedRemote: "172.16.0.1",
		},
		{
			name:           "No headers - use original RemoteAddr",
			headers:        map[string]string{},
			originalRemote: "127.0.0.1:8080",
			expectedRemote: "127.0.0.1:8080",
		},
		{
			name: "Empty header values",
			headers: map[string]string{
				"True-Client-IP":  "",
				"X-Real-IP":       "",
				"X-Forwarded-For": "",
			},
			originalRemote: "127.0.0.1:8080",
			expectedRemote: "127.0.0.1:8080",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = test.originalRemote

			// Set headers
			for key, value := range test.headers {
				req.Header.Set(key, value)
			}

			// Create response recorder
			recorder := httptest.NewRecorder()

			// Execute middleware
			wrappedHandler.ServeHTTP(recorder, req)

			// Check result
			if capturedRemoteAddr != test.expectedRemote {
				t.Errorf("Expected RemoteAddr %q, got %q", test.expectedRemote, capturedRemoteAddr)
			}
		})
	}
}

func TestRequestID(t *testing.T) {
	// Create middleware
	middlewareFunc := RequestID()
	if middlewareFunc == nil {
		t.Fatal("RequestID() returned nil middleware")
	}

	// Test handler that captures request ID from context
	var capturedRequestID string
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get request ID from context using chi's key
		capturedRequestID = middleware.GetReqID(r.Context())
	})

	// Wrap test handler with middleware
	wrappedHandler := middlewareFunc(testHandler)

	t.Run("Injects request ID into context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(recorder, req)

		if capturedRequestID == "" {
			t.Error("Request ID should be injected into context")
		}
	})

	t.Run("Request ID has expected format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(recorder, req)

		if capturedRequestID == "" {
			t.Fatal("Request ID not found")
		}

		// Should contain a slash separating hostname/process from counter
		if !strings.Contains(capturedRequestID, "/") {
			t.Errorf("Request ID should contain '/', got: %s", capturedRequestID)
		}

		// Should end with a number (counter)
		parts := strings.Split(capturedRequestID, "/")
		if len(parts) < 2 {
			t.Errorf("Request ID should have format 'host/counter', got: %s", capturedRequestID)
		}
	})

	t.Run("Multiple requests get different IDs", func(t *testing.T) {
		var requestIDs []string

		// Make multiple requests
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest("GET", "/", nil)
			recorder := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(recorder, req)
			requestIDs = append(requestIDs, capturedRequestID)
		}

		// All request IDs should be different
		for i := 0; i < len(requestIDs); i++ {
			for j := i + 1; j < len(requestIDs); j++ {
				if requestIDs[i] == requestIDs[j] {
					t.Errorf("Request IDs should be unique, but got duplicate: %s", requestIDs[i])
				}
			}
		}
	})
}

func TestRequestID_WithoutMiddleware(t *testing.T) {
	// Test that without middleware, no request ID is in context
	var capturedRequestID string
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequestID = middleware.GetReqID(r.Context())
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	testHandler.ServeHTTP(recorder, req)

	if capturedRequestID != "" {
		t.Errorf("Without middleware, request ID should be empty, got: %s", capturedRequestID)
	}
}

func TestRecoverer(t *testing.T) {
	// Create middleware
	middlewareFunc := Recoverer()
	if middlewareFunc == nil {
		t.Fatal("Recoverer() returned nil middleware")
	}

	t.Run("Recovers from panic and returns 500", func(t *testing.T) {
		// Handler that panics
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})

		// Wrap with recoverer middleware
		wrappedHandler := middlewareFunc(panicHandler)

		req := httptest.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		// This should not panic the test
		wrappedHandler.ServeHTTP(recorder, req)

		// Should return 500 status
		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", recorder.Code)
		}
	})

	t.Run("Normal requests work fine", func(t *testing.T) {
		// Normal handler that doesn't panic
		normalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		wrappedHandler := middlewareFunc(normalHandler)

		req := httptest.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(recorder, req)

		// Should return 200 status
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", recorder.Code)
		}

		// Should have response body
		if recorder.Body.String() != "OK" {
			t.Errorf("Expected body 'OK', got %q", recorder.Body.String())
		}
	})

	t.Run("Recovers from panic with request ID", func(t *testing.T) {
		// Handler that panics
		panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("test panic with request ID")
		})

		// Wrap with both RequestID and Recoverer middleware
		requestIDMiddleware := RequestID()
		recovererMiddleware := Recoverer()

		// Chain middlewares: RequestID first, then Recoverer
		wrappedHandler := recovererMiddleware(requestIDMiddleware(panicHandler))

		req := httptest.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		// This should not panic the test
		wrappedHandler.ServeHTTP(recorder, req)

		// Should return 500 status
		if recorder.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", recorder.Code)
		}
	})

	t.Run("Recovers from different panic types", func(t *testing.T) {
		testCases := []struct {
			name      string
			panicFunc func()
		}{
			{
				name:      "string panic",
				panicFunc: func() { panic("string panic") },
			},
			{
				name:      "error panic",
				panicFunc: func() { panic(http.ErrHandlerTimeout) },
			},
			{
				name:      "nil panic",
				panicFunc: func() { panic(nil) },
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					tc.panicFunc()
				})

				wrappedHandler := middlewareFunc(panicHandler)

				req := httptest.NewRequest("GET", "/", nil)
				recorder := httptest.NewRecorder()

				// Should not panic the test
				wrappedHandler.ServeHTTP(recorder, req)

				// Should return 500 status for all panic types
				if recorder.Code != http.StatusInternalServerError {
					t.Errorf("Expected status 500 for %s, got %d", tc.name, recorder.Code)
				}
			})
		}
	})
}

func TestRecoverer_WithoutMiddleware(t *testing.T) {
	// Test that without middleware, panics are not recovered
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unrecovered panic")
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	// This should panic and be caught by the test framework
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic to occur without recoverer middleware")
		}
	}()

	panicHandler.ServeHTTP(recorder, req)
}
