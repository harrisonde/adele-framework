package main

import (
	"errors"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/harrisonde/adel"
)

var MakeKeyCommand = &adel.Command{
	Name:        "make key",
	Help:        "create new application key",
	Description: "use the make key command to create a new application key for your adel application",
	Usage:       "make session",
	Options: map[string]string{
		"-f, --force":    "force replace application key, even if it exists",
		"-g, --generate": "generate a key and print to cli",
	},
}

func doMakeKey() error {
	color.Yellow("Starting key generation")
	color.Green("  Creating new application key...")
	rnd := ade.RandomString(32)
	longOption, _ := GetOption("generate")
	shortOption, _ := GetOption("g")
	if longOption == "generate" || shortOption == "g" {
		color.Green("  found option to generate application key and print to cli...")
		color.Yellow("Encryption key:")
		color.White("  " + rnd)
	} else {
		path := ade.RootPath + "/.env"
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		r := regexp.MustCompile("APP_KEY=.*")
		hasKey := r.FindStringSubmatch(string(read))

		if len(hasKey) < 1 {
			return errors.New("unable to find the application key in the env file")
		}

		if len(hasKey[0]) > 8 {
			longOption, _ := GetOption("force")
			shortOption, _ := GetOption("f")
			if longOption == "force" || shortOption == "f" {
				color.Green("  option to force application key replacement found")
			} else {
				return errors.New("the application key exists; use -f, --force to replace")
			}
		}

		color.Green("  Write application key to file...")
		newEnv := r.ReplaceAll(read, []byte("APP_KEY="+rnd))
		err = os.WriteFile(ade.RootPath+"/.env", []byte(newEnv), 0)
		if err != nil {
			return err
		}
		color.Yellow("Application key generation complete")

	}

	return nil
}
