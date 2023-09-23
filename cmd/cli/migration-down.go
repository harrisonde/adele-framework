package main

import "github.com/harrisonde/adel"

var MakeMigrateDownCommand = &adel.Command{
	Name: "migrate down",
	Help: "reverse the most recent migration",
}
