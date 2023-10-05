package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/harrisonde/adel"
)

var ServerCommand = &adel.Command{
	Name:        "server",
	Help:        "manage the application server",
	Description: "use the server command to start or place server in maintenance mode",
	Usage:       "serve [options]",
	Options: map[string]string{
		"-d, --down": "take the server out of maintenance mode and handel requests",
		"-u, --up":   "put the server in maintenance mode",
	},
}

func handelServer() error {

	hasLongFormat := HasOption("down")
	hasShortFormat := HasOption("d")
	if hasLongFormat || hasShortFormat {
		setMaintenanceMode(true)
		return nil
	}

	hasLongFormat = HasOption("up")
	hasShortFormat = HasOption("u")
	if hasLongFormat || hasShortFormat {
		setMaintenanceMode(false)
		return nil
	}

	err := boot()
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func boot() error {
	r := isStarted()
	if r == true {
		return errors.New("application server is already running")
	}

	build()
	_, err := os.Stat(ade.RootPath + "/tmp/" + ade.AppName)
	if os.IsNotExist(err) {
		writeOutput("Your application binary was not found. Did you try building it before starting the server? If not, please build your application binary by running: \n\t$ go build\n")
		return errors.New("Cannot start application server without a binary.")
	}

	writeOutput("Adel is running on "+os.Getenv("APP_URL"), "info")
	writeOutput("Press Ctrl+C to stop the server\n")

	cmd := exec.Command(ade.RootPath + "/tmp/" + ade.AppName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func build() {
	cmd := exec.Command("go", "build", "-o", ade.RootPath+"/tmp/"+ade.AppName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		exitGracefully(err)
	}
}

func setMaintenanceMode(inMaintenanceMode bool) {

	client, err := getRpcClient()
	if err != nil {
		exitGracefully(err)
	}

	var result string
	err = client.Call("RPCServer.MaintenanceMode", inMaintenanceMode, &result)
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow(result)
}

func isStarted() bool {
	_, err := getRpcClient()
	return err == nil
}
