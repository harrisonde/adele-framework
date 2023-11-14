package main

import "github.com/harrisonde/adele-framework"

var MakeMigrateUpCommand = &adele.Command{
	Name:        "migrate up",
	Help:        "run all up migrations",
	Description: "use the migrate up command to run all migrations",
	Usage:       "migrate up",
	Options:     map[string]string{},
}

func doMigrateUp() error {

	tx, err := ade.PopConnect()
	if err != nil {
		exitGracefully(err)
	}

	defer tx.Close()

	err = ade.RunPopMigrations(tx)
	if err != nil {
		return err
	}

	return nil
}
