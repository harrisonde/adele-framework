package main

import "github.com/harrisonde/adele"

var MakeMigrateResetCommand = &adele.Command{
	Name:        "migrate reset",
	Help:        "reset and re-run all migrations",
	Description: "use the migrate reset command to run all the down migrations and all the up migrations.",
	Usage:       "migrate up",
	Options:     map[string]string{},
}

func doMigrateReset() error {
	tx, err := ade.PopConnect()
	if err != nil {
		exitGracefully(err)
	}

	defer tx.Close()

	err = ade.PopMigrateReset(tx)
	if err != nil {
		return err
	}

	return nil
}
