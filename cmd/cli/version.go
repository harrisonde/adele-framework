package main

import (
	"github.com/fatih/color"
	"github.com/harrisonde/adele-framework"
)

var VersionCommand = &adele.Command{
	Name: "version",
	Help: "print application version",
}

const version = "1.0.0"

func printVersion() {
	color.Yellow("Application version: " + version)
}
