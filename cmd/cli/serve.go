package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/harrisonde/adel"
)

var ServerCommand = &adel.Command{
	Name: "serve",
	Help: "start the application server to handle http requests",
}

func doStart() error {

	build()
	_, err := os.Stat(ade.RootPath + "/tmp/" + ade.AppName)
	if os.IsNotExist(err) {
		writeOutput("Your application binary was not found. Did you try building it before starting the server? If not, please build your application binary by running: \n\t$ go build\n")
		exitGracefully(errors.New("Cannot start application server without a binary."))
	}

	writeOutput("Adel is running on "+os.Getenv("APP_URL"), "info")
	writeOutput("Press Ctrl+C to stop the server\n")

	cmd := exec.Command(ade.RootPath + "/tmp/" + ade.AppName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitGracefully(err)
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
