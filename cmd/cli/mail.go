package main

import "github.com/harrisonde/adel"

var MakeMailCommand = &adel.Command{
	Name: "make mail <name>",
	Help: "create two stub mail templates in the mail directory",
}
