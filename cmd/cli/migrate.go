package main

import "github.com/harrisonde/adele-framework"

var MakeMigrateCommand = &adele.Command{
	Name:        "migrate",
	Help:        "run all migrations",
	Description: "use the migrate command to run all migrations that have not be previously executed",
	Usage:       "migrate",
	Options:     map[string]string{},
}

func doMigrate(arg2, arg3 string) error {

	checkForDb()

	switch arg2 {
	case "up":
		err := doMigrateUp()
		if err != nil {
			return err
		}

	case "down":
		err := doMigrateDown()
		if err != nil {
			return err
		}
	case "reset":
		err := doMigrateReset()
		if err != nil {
			return err
		}

	default:
		showHelp()
	}

	return nil
}
