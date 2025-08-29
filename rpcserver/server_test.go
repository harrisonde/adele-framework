package rpcserver

import (
	"os"
	"testing"

	"github.com/cidekar/adele-framework"
)

func TestServerStart_InvalidPort(t *testing.T) {
	// Make sure RPC is not disabled
	os.Unsetenv("RPC_SERVER_DISABLE")

	// Set an invalid port for testing
	os.Setenv("RPC_SERVER_PORT", "999999999") // Invalid port that should fail to bind

	defer os.Unsetenv("RPC_SERVER_PORT")

	app := &adele.Adele{}

	err := Start(app)
	// This should fail because port 99999 is likely not available or invalid
	if err == nil {
		t.Fatal("Start() should have failed with invalid/unavailable port")
	}
}

func TestServerStop_RPCDisabled(t *testing.T) {
	// Test that Stop returns nil when RPC is disabled
	os.Setenv("RPC_SERVER_DISABLE", "true")
	defer os.Unsetenv("RPC_SERVER_DISABLE")

	app := &adele.Adele{}
	err := Stop(app)
	if err != nil {
		t.Fatalf("Stop() should return nil when RPC_SERVER_DISABLE is set: %v", err)
	}
}

func TestServerStop_NilApplication(t *testing.T) {
	// Test Stop with nil application - should not panic
	os.Unsetenv("RPC_SERVER_DISABLE")

	err := Stop(nil)
	// This might return an error or panic - we just want to make sure it's handled
	_ = err // Don't assert on error since we don't know exact behavior
}

func TestRPCServer_SetMaintenanceMode(t *testing.T) {
	// Create RPCServer with a simple Application that has an Adele instance
	app := &adele.Adele{}
	server := &RPCServer{App: app}

	// Test setting to true
	args := &MaintenanceModeArgs{InMaintenanceMode: true}
	reply := &MaintenanceModeReply{}

	err := server.SetMaintenanceMode(args, reply)
	if err != nil {
		t.Fatalf("SetMaintenanceMode(true) failed: %v", err)
	}
	if !app.MaintenanceMode {
		t.Error("MaintenanceMode should be true")
	}
	if reply.Status != "down" {
		t.Errorf("Expected status 'down', got '%s'", reply.Status)
	}

	// Test setting to false
	args.InMaintenanceMode = false
	reply.Status = "" // Reset

	err = server.SetMaintenanceMode(args, reply)
	if err != nil {
		t.Fatalf("SetMaintenanceMode(false) failed: %v", err)
	}
	if app.MaintenanceMode {
		t.Error("MaintenanceMode should be false")
	}
	if reply.Status != "up" {
		t.Errorf("Expected status 'up', got '%s'", reply.Status)
	}
}
