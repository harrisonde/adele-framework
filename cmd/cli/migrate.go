package main

func doMigrate(arg2, arg3 string) error {

	checkForDb()

	tx, err := ade.PopConnect()
	if err != nil {
		exitGracefully(err)
	}

	defer tx.Close()

	switch arg2 {
	case "up":
		err := ade.RunPopMigrations(tx)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
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

	case "reset":
		err := ade.PopMigrateReset(tx)
		if err != nil {
			return err
		}

	default:
		showHelp()
	}

	return nil
}
