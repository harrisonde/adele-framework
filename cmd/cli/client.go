package main

import "github.com/harrisonde/adel"

var MakeClientCommand = &adel.Command{
	Name: "make client <name>",
	Help: "makes a oauth2 password grant client",
}
