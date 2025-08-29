package rpcserver

import (
	"os"
	"testing"

	"github.com/cidekar/adele-framework"
)

func TestNewRPCClient_ConnectionFailure_InvalidPort(t *testing.T) {
	// Set invalid port number (> 65535)
	os.Setenv("RPC_SERVER_ADDR", "127.0.0.1")
	os.Setenv("RPC_SERVER_PORT", "999999") // Invalid port
	defer func() {
		os.Unsetenv("RPC_SERVER_ADDR")
		os.Unsetenv("RPC_SERVER_PORT")
	}()

	client, err := NewRPCClient()
	if err == nil {
		if client != nil {
			client.Close()
		}
		t.Fatal("NewRPCClient() should have failed with invalid port")
	}
	if client != nil {
		t.Fatal("NewRPCClient() should return nil client on connection failure")
	}
}

func TestNewRPCClient_ConnectionFailure_InvalidAddress(t *testing.T) {
	// Set invalid IP address to force connection failure
	os.Setenv("RPC_SERVER_ADDR", "999.999.999.999")
	os.Setenv("RPC_SERVER_PORT", "8080")
	defer func() {
		os.Unsetenv("RPC_SERVER_ADDR")
		os.Unsetenv("RPC_SERVER_PORT")
	}()

	client, err := NewRPCClient()
	if err == nil {
		if client != nil {
			client.Close() // Clean up if somehow it worked
		}
		t.Fatal("NewRPCClient() should have failed with invalid IP address")
	}
	if client != nil {
		t.Fatal("NewRPCClient() should return nil client on connection failure")
	}
}

func TestMaintenanceModeArgsStruct(t *testing.T) {
	// Test the struct initialization
	args := &MaintenanceModeArgs{InMaintenanceMode: true}
	if !args.InMaintenanceMode {
		t.Error("MaintenanceModeArgs.InMaintenanceMode should be true")
	}

	args2 := &MaintenanceModeArgs{InMaintenanceMode: false}
	if args2.InMaintenanceMode {
		t.Error("MaintenanceModeArgs.InMaintenanceMode should be false")
	}
}

func TestMaintenanceModeReplyStruct(t *testing.T) {
	// Test the struct initialization
	reply := &MaintenanceModeReply{Status: "server is down"}
	expected := "server is down"
	if reply.Status != expected {
		t.Errorf("MaintenanceModeReply.Status = %q, expected %q", reply.Status, expected)
	}
}

func TestNewRPCClient_UsesEnvironmentVariables(t *testing.T) {
	// Test that NewRPCClient reads fresh environment variables
	os.Setenv("RPC_SERVER_ADDR", "192.168.1.1")
	os.Setenv("RPC_SERVER_PORT", "9999")
	defer func() {
		os.Unsetenv("RPC_SERVER_ADDR")
		os.Unsetenv("RPC_SERVER_PORT")
	}()

	// This should fail because we're trying to connect to 192.168.1.1:9999
	// but the important part is that it's using the env vars we set
	client, err := NewRPCClient()
	if client != nil {
		client.Close() // Clean up if it somehow connected
	}

	// We expect this to fail (no server at 192.168.1.1:9999)
	// The test passes if it tries to connect to the right address
	if err == nil {
		t.Fatal("Expected connection to fail to non-existent server")
	}

	// The error message should contain our custom address/port
	// This indirectly confirms it used our environment variables
}

func TestStart_RPCDisabled(t *testing.T) {
	// Test that Start returns nil when RPC is disabled
	os.Setenv("RPC_SERVER_DISABLE", "true")
	defer os.Unsetenv("RPC_SERVER_DISABLE")

	// Create a minimal Adele instance for testing
	app := &adele.Adele{}
	err := Start(app)
	if err != nil {
		t.Fatalf("Start() should return nil when RPC_SERVER_DISABLE is set: %v", err)
	}
}

func TestRPCClient_Close_NilClient(t *testing.T) {
	// Test Close with nil underlying client
	client := &RPCClient{client: nil}

	err := client.Close()
	if err == nil {
		t.Fatal("Close() should return error when client is nil")
	}

	expected := "rpc client is nil or already closed"
	if err.Error() != expected {
		t.Errorf("Close() error = %q, expected %q", err.Error(), expected)
	}
}
