package main

import (
	"github.com/fatih/color"
	"github.com/harrisonde/adele-framework"
)

var VersionCommand = &adele.Command{
	Name: "version",
	Help: "print application version",
}

func printVersion() {
	color.Yellow("Application version: " + ade.Version)
}
