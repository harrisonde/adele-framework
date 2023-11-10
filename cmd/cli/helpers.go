package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/harrisonde/adele/cmd"
	"github.com/joho/godotenv"
)

var optionPattern = "(^--[\\w\\d]{0,}|^-[\\w\\d]{0,})"

func showHelp(arg1 ...string) {
	var help string

	if len(arg1) > 0 {
		help = cmd.GetHelp(arg1[0])
	} else {
		help = cmd.GetHelp()
	}

	color.Yellow(help)
}

func loadOptions() {
	args := os.Args[1:]
	shift := 0
	for k, v := range args {
		isOpt, _ := regexp.MatchString(optionPattern, v)
		if isOpt {
			cmdOptions = append(cmdOptions, v)
			os.Args[k+1] = ""
			shift++
		} else if shift > 0 {
			shiftToPos := k - shift + 1
			os.Args[shiftToPos] = v
		}
	}
}
func setup(arg1, arg2, arg3 string) {
	if HasOption("help") || HasOption("h") {
		if arg1 == "" {
			showHelp()
		} else if arg2 != "" {
			subCmd := arg1 + " " + arg2
			showHelp(subCmd)
		} else {
			showHelp(arg1)
		}
		os.Exit(0)
	}

	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		err := godotenv.Load()
		if err != nil {
			exitGracefully(err)
		}

		path, err := os.Getwd()
		if err != nil {
			exitGracefully(err)
		}

		ade.RootPath = path
		ade.AppName = os.Getenv("APP_NAME")
		ade.DB.DataType = os.Getenv("DATABASE_TYPE")
	}
}

func getDSN() string {
	dbType := ade.DB.DataType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASSWORD") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASSWORD"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		return dsn
	}

	// mariadb / sql
	return "mysql://" + ade.BuildDSN()
}

func checkForDb() {
	dbType := ade.DB.DataType

	if dbType == "" {
		exitGracefully(errors.New("no database connection provided in .env. Did you create one?"))
	}

	if !fileExists(ade.RootPath + "/config/database.yml") {
		exitGracefully(errors.New("config/database.yml does not exist. Did you create one?"))
	}
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	matched, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}
	if matched {
		read, err := os.ReadFile(path)
		if err != nil {
			exitGracefully(err)
		}
		newCont := strings.Replace(string(read), "myapp", appURL, -1)
		err = os.WriteFile(path, []byte(newCont), 0)
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}

func updateSource() {
	err := filepath.Walk(".", updateSourceFiles)
	if err != nil {
		exitGracefully(err)
	}
}

func GetOption(option string) (string, error) {
	var o string
	for _, v := range cmdOptions {
		oc := strings.ReplaceAll(v, "-", "")

		// --switch will take no arguments
		// --switch= option accepts argument
		// --switch=foo option accepts argument, default foo

		switchAcceptArg := strings.Split(oc, "=")
		if len(switchAcceptArg) == 1 {
			if oc == option {
				o = oc
			}
		} else if len(switchAcceptArg) == 2 {
			o = switchAcceptArg[1]
		} else {
			if oc == option {
				o = oc
			}
		}
	}

	if o == "" {
		return "", errors.New(fmt.Sprint("%s is not a known argument", option))
	}
	return o, nil
}

func HasOption(option string) bool {
	for _, v := range cmdOptions {
		oc := strings.ReplaceAll(v, "-", "")
		switchAcceptArg := strings.Split(oc, "=")
		if len(switchAcceptArg) == 1 {
			if switchAcceptArg[0] == option {
				return true
			}
		} else if len(switchAcceptArg) == 2 {
			if switchAcceptArg[0] == option {
				return true
			}
		}
	}
	return false
}
