package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/harrisonde/adel"
)

const version = "1.0.0"

var ade adel.Adel
var message string

func main() {
	arg1, arg2, arg3, arg4, err := validateInput()
	if err != nil {
		exitGracefully(err)
	}

	setup(arg1, arg2)

	switch arg1 {
	case "help":
		showHelp()

	case "up":
		rpcClient(false)

	case "down":
		rpcClient(true)

	case "serve":
		doStart()

	case "new":
		if arg2 == "" {
			exitGracefully(errors.New("new requires an application name"))
		}
		doNew(arg2)

	case "version":
		color.Yellow("Application version: " + version)

	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		err = doMigrate(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
		message = "Migrations complete."

	case "make":
		if arg2 == "" {
			exitGracefully(errors.New("make requires a subcommand: (migration|model|handler)"))
		}

		err = doMake(arg2, arg3, arg4)

		if err != nil {
			exitGracefully(err)
		}

	default:
		showHelp()
	}

	exitGracefully(nil, message)
}

func validateInput() (string, string, string, string, error) {
	var arg1, arg2, arg3, arg4 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}

		if len(os.Args) >= 5 {
			arg4 = os.Args[4]
		}

	} else {
		color.Red("Error: command required")
		showHelp()
		return "", "", "", "", errors.New("command required")
	}
	return arg1, arg2, arg3, arg4, nil
}

func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Finished!")
	}

	os.Exit(0)
}

func writeOutput(message string, writeType ...string) {
	w := ""

	if len(writeType) > 0 {
		w = writeType[0]
	}

	switch w {
	case "info":
		notice := color.New(color.BgBlue, color.FgWhite).SprintFunc()
		fmt.Printf("\n %s %s\n", notice(" INFO "), message)

	case "line":
	default:
		fmt.Printf("\n %s\n", message)
	}

}
