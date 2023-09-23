package main

import "github.com/harrisonde/adel"

var MakeHandlerCommand = &adel.Command{
	Name: "make handler <name>",
	Help: "create a stub handler in the handlers directory",
}
