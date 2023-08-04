package main

import (
	"os"
	"os/exec"
)

func doStart() error {

	cmd := exec.Command("go", "mod", "vendor")
	err := cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	cmd = exec.Command("go", "build", "-o", ade.RootPath+"/tmp/adelApp")
	err = cmd.Start()
	if err != nil {
		exitGracefully(err)
	}

	writeOutput("Adel is running on "+os.Getenv("APP_URL"), "info")
	writeOutput("Press Ctrl+C to stop the server\n")

	cmd = exec.Command(ade.RootPath + "/tmp/adelApp")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitGracefully(err)
	}

	return nil

}
