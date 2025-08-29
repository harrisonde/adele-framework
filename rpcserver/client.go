package rpcserver

import (
	"errors"
	"net/rpc"

	"github.com/cidekar/adele-framework"
)

type MaintenanceModeArgs struct {
	InMaintenanceMode bool
}

type MaintenanceModeReply struct {
	Status string
}

type RPCClient struct {
	client *rpc.Client
}

// Creates a new RPC connection to a server and wraps the connection in
// a custom RPCClient struct. The client is returned ready to use, or
// an error if connection fails.
func NewRPCClient() (*RPCClient, error) {
	ServerAddr := adele.Helpers.Getenv("RPC_SERVER_ADDR", ServerAddrDefault)
	ServerPort := adele.Helpers.Getenv("RPC_SERVER_PORT", ServerPortDefault)
	client, err := rpc.Dial("tcp", ServerAddr+":"+ServerPort)
	if err != nil {
		return nil, err
	}

	return &RPCClient{client: client}, nil
}

// Closes the underlying RPC connection using a clean shutdown pattern. Any
// errors duing the close are passed through from the underlying close
// operation.
func (c *RPCClient) Close() error {
	if c.client == nil {
		return errors.New("rpc client is nil or already closed")
	}
	return c.client.Close()
}

// Convenience wrapper method that makes RPC calls easier and more user-friendly— the
// the method is used to put the server in or out of maintance mode.
// Example usage:
//
//	status, err := client.SetMaintenanceMode(true)  // Put in maintenance
//	status, err := client.SetMaintenanceMode(false) // Take out of maintenance
func (c *RPCClient) SetMaintenanceMode(inMaintenance bool) (string, error) {
	args := &MaintenanceModeArgs{InMaintenanceMode: inMaintenance}
	reply := &MaintenanceModeReply{}

	err := c.client.Call("RPCServer.SetMaintenanceMode", args, reply)
	return reply.Status, err
}
