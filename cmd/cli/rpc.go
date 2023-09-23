package main

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/fatih/color"
	"github.com/harrisonde/adel"
)

var RpcCommandUp = &adel.Command{
	Name: "up",
	Help: "take the server out of maintenance mode",
}

var RpcCommandDown = &adel.Command{
	Name: "down",
	Help: "put the server in maintenance mode",
}

func rpcClient(inMaintenanceMode bool) {
	port := os.Getenv("RPC_PORT")
	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		exitGracefully(err)
	}

	fmt.Println("Connected...")

	var result string
	err = client.Call("RPCServer.MaintenanceMode", inMaintenanceMode, &result)
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow(result)
}

// May need to pass the commands from the app here as an argument
func rpcCommand(command string) {
	port := os.Getenv("RPC_PORT")
	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		exitGracefully(err)
	}

	var result string
	err = client.Call("RPCServer.Command", command, &result)
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow(result)
}
