package main

import "github.com/harrisonde/adele"

var MakeMigrateDownCommand = &adele.Command{
	Name:        "migrate down",
	Help:        "reverse the last migration",
	Description: "use the migrate down command to reverse the most recent migration",
	Usage:       "migrate down",
	Options: map[string]string{
		"-a, --all": "reverse all migrations",
	},
}

func doMigrateDown() error {

	tx, err := ade.PopConnect()
	if err != nil {
		exitGracefully(err)
	}

	defer tx.Close()

	longOption, _ := GetOption("skip")
	shortOption, _ := GetOption("s")
	if longOption == "all" || shortOption == "a" {
		err := ade.PopMigrateDown(tx, -1)
		if err != nil {
			return err
		}
	} else {
		err := ade.PopMigrateDown(tx, 1)
		if err != nil {
			return err
		}
	}
	return nil
}
