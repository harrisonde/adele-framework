package main

import "github.com/harrisonde/adel"

var MakeMigrateUpCommand = &adel.Command{
	Name: "migrate up",
	Help: "run all up migrations",
}
