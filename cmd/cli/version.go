package main

import (
	"github.com/fatih/color"
	"github.com/harrisonde/adel"
)

var VersionCommand = &adel.Command{
	Name: "version",
	Help: "print application version",
}

const version = "1.0.0"

func printVersion() {
	color.Yellow("Application version: " + version)
}
