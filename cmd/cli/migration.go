package main

import (
	"errors"

	"github.com/harrisonde/adel"
)

var MakeMigrationCommand = &adel.Command{
	Name:        "make migration",
	Help:        "create a new migration",
	Description: "use the make migration command to create a new migration template in the migrations directory",
	Usage:       "migrate --type=fizz",
	Options: map[string]string{
		"-t=fizz, --type=fizz": "type of migration to create sql/fizz (default is --type=fizz)",
	},
}

func doMakeMigration(arg3, arg4 string) error {
	checkForDb()

	if arg3 == "" {
		return errors.New("You must provide a name for the migration!")
	}

	var up, down string

	migrationType := "fizz"
	supported := []string{
		"fizz",
		"sql",
	}

	hasLongFormat := HasOption("type")
	hasShortFormat := HasOption("t")
	if hasLongFormat || hasShortFormat {
		mt, _ := GetOption("type")
		if mt == "" {
			mt, _ = GetOption("t")
		}
		for _, s := range supported {
			if mt == s {
				migrationType = s
			}
		}
	}

	if migrationType == "fizz" {
		upBytes, err := templateFS.ReadFile("templates/migrations/migration_up.fizz")
		if err != nil {
			return err
		}

		downBytes, err := templateFS.ReadFile("templates/migrations/migration_down.fizz")
		if err != nil {
			return err
		}

		up = string(upBytes)
		down = string(downBytes)
	} else {
		migrationType = "sql"
	}

	err := ade.CreatePopMigration([]byte(up), []byte(down), arg3, migrationType)
	if err != nil {
		return err
	}

	return nil
}
