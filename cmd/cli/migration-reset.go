package main

import "github.com/harrisonde/adel"

var MakeMigrateResetCommand = &adel.Command{
	Name: "migrate reset",
	Help: "run all down migrations and all up migrations",
}
