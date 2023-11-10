package main

import (
	"github.com/harrisonde/adele"
	"github.com/harrisonde/adele/cmd"
)

var MakeCommand = &adele.Command{
	Name:        "make",
	Help:        "show all make commands",
	Description: "use the make command to create a new resources e.g., migrations, models, or mail",
	Usage:       "make <command> [options] [arguments]",
	Options:     map[string]string{},
}

func doMake(arg2, arg3, arg4 string) error {

	switch arg2 {

	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}
		return nil

	case "oauth":
		err := doOauth()
		if err != nil {
			exitGracefully(err)
		}
		return nil

	case "mail":
		err := doMakeMail(arg3)
		if err != nil {
			exitGracefully(err)
		}
		return nil

	case "migration":
		err := doMakeMigration(arg3, arg4)
		if err != nil {
			exitGracefully(err)
		}
		return nil

	case "model":
		err := doMakeModel(arg3)
		if err != nil {
			exitGracefully(err)
		}
		return nil

	case "handler":
		err := doMakeHandler(arg3)
		if err != nil {
			exitGracefully(err)
		}
		return nil

	case "key":
		err := doMakeKey()
		if err != nil {
			exitGracefully(err)
		}
		return nil

	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}
		return nil
	}

	cmd.GetHelp("make")
	return nil
}
