package main

import (
	"errors"
	"os"
	"os/exec"
)

func doStart() error {
	_, err := os.Stat(ade.RootPath + "/" + ade.AppName)
	if os.IsNotExist(err) {
		writeOutput("Your application binary was not found. Did you try building it before starting the server? If not, please build your application binary by running: \n\t$ go build\n")
		exitGracefully(errors.New("Cannot start application server without a binary."))
	}

	writeOutput("Adel is running on "+os.Getenv("APP_URL"), "info")
	writeOutput("Press Ctrl+C to stop the server\n")

	cmd := exec.Command(ade.RootPath + "/" + ade.AppName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitGracefully(err)
	}

	return nil
}
