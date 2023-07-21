package main

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/fatih/color"
)

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
