package httpserver

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cidekar/adele-framework"
	"github.com/cidekar/adele-framework/mux"
	"github.com/sirupsen/logrus"
)

func TestNewServer_DefaultPort(t *testing.T) {
	// Clear PORT environment variable
	os.Unsetenv("PORT")

	// Create minimal Adele instance
	app := &adele.Adele{
		Routes: mux.NewRouter(),
		Log:    logrus.New(),
	}

	server := NewServer(app)
	if server.Addr != ":4000" {
		t.Errorf("Expected default addr ':4000', got %s", server.Addr)
	}
}

func TestNewServer_CustomPort(t *testing.T) {
	// Set custom port
	os.Setenv("HTTP_PORT", "8080")
	defer os.Unsetenv("PORT")

	// Create minimal Adele instance
	app := &adele.Adele{
		Routes: mux.NewRouter(),
		Log:    logrus.New(),
	}

	server := NewServer(app)
	if server.Addr != ":8080" {
		t.Errorf("Expected addr ':8080', got %s", server.Addr)
	}
}

func TestNewServer_Configuration(t *testing.T) {
	// Test all server configuration
	app := &adele.Adele{
		Routes: mux.NewRouter(),
		Log:    logrus.New(),
	}

	server := NewServer(app)

	// Test timeouts
	if server.IdleTimeout != 30*time.Second {
		t.Errorf("Expected IdleTimeout 30s, got %v", server.IdleTimeout)
	}
	if server.ReadTimeout != 30*time.Second {
		t.Errorf("Expected ReadTimeout 30s, got %v", server.ReadTimeout)
	}
	if server.WriteTimeout != 600*time.Second {
		t.Errorf("Expected WriteTimeout 600s, got %v", server.WriteTimeout)
	}

	// Test handler is set
	if server.Handler == nil {
		t.Error("Handler should not be nil")
	}

	// Test ErrorLog is set
	if server.ErrorLog == nil {
		t.Error("ErrorLog should not be nil")
	}
}

func TestNewServer_NilAdele(t *testing.T) {
	// Test with nil adele - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil adele to NewServer")
		}
	}()

	NewServer(nil)
}

func TestStart_ImmediateError(t *testing.T) {
	// Use port 22 which is ssh and should be busy
	os.Setenv("HTTP_PORT", "22")
	defer os.Unsetenv("HTTP_PORT")

	app := &adele.Adele{
		Routes: mux.NewRouter(),
		Log:    logrus.New(),
	}

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		err := Start(app)
		errChan <- err
	}()

	// Wait for error with timeout
	select {
	case err := <-errChan:
		if err == nil {
			t.Fatal("Expected error when trying to bind to port 1 (requires root)")
		}
		// Error should mention permission or bind issue
		errorStr := strings.ToLower(err.Error())
		if !strings.Contains(errorStr, "permission") && !strings.Contains(errorStr, "bind") && !strings.Contains(errorStr, "address") {
			t.Errorf("Expected error about permission/bind/address, got: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out - server should have failed immediately")
	}
}

func TestStart_NilAdele(t *testing.T) {
	// Test with nil adele - should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil adele to Start")
		}
	}()

	Start(nil)
}

func TestStart_InvalidPort(t *testing.T) {
	// Set invalid port
	os.Setenv("HTTP_PORT", "invalid")
	defer os.Unsetenv("HTTP_PORT")

	app := &adele.Adele{
		Routes: mux.NewRouter(),
		Log:    logrus.New(),
	}

	err := Start(app)
	if err == nil {
		t.Fatal("Expected error with invalid port")
	}

	// Should get an error about invalid port/address
	if !strings.Contains(err.Error(), "invalid") && !strings.Contains(err.Error(), "address") {
		t.Errorf("Expected error about invalid port/address, got: %v", err)
	}
}
