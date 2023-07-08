package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	switch arg2 {
	case "up":
		err := ade.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
			err := ade.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := ade.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := ade.MigrateDownAll(dsn)
		if err != nil {
			return err
		}

		err = ade.MigrateUp(dsn)
		if err != nil {
			return err
		}

	default:
		showHelp()
	}

	return nil
}
