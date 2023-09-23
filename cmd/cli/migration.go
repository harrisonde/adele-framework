package main

import "github.com/harrisonde/adel"

var MakeMigrationCommand = &adel.Command{
	Name: "make migration <name> <format>",
	Help: "create a new migration; format=sql/fizz (default fizz)",
}
