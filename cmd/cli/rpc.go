package main

import (
	"net/rpc"
	"os"
	"strings"

	"github.com/fatih/color"
)

func getRpcClient() (*rpc.Client, error) {
	port := os.Getenv("RPC_PORT")
	client, err := rpc.Dial("tcp", "127.0.0.1:"+port)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func rpcCommand(arg1, arg2, arg3 string) {

	client, err := getRpcClient()
	if err != nil {
		showHelp()
		exitGracefully(err)
	}

	var result string
	options := strings.Join(cmdOptions, ",")
	command := []string{
		arg1, arg2, arg3, options,
	}

	err = client.Call("RPCServer.Command", command, &result)
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow(result)
}
