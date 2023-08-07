package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup(arg1, arg2 string) {
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

func showHelp() {
	color.Yellow(`Available commands:

	help                           - show help commands
	up                             - take the server out of maintenance mode
	down                           - put the server in maintenance mode
	serve                          - start the application server to handle http requests
	make                           - show all make commands
	make auth                      - install authentication
	make client <name>             - makes a oauth2 password grant client
	make handler <name>            - create a stub handler in the handlers directory
	make mail <name>               - create two stub mail templates in the mail directory
	make migration <name> <format> - create a new migration; format=sql/fizz (default fizz)
	make model <name>              - create a new model in the data directory
	make session                   - create a table in the database to store sessions
	migrate                        - run all migration that have not been run
	migrate down                   - reverse the most recent migration
	migrate reset                  - run all down migrations and all up migrations
	version                        - print application version

	`)
}

func updateSourceFiles(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// is a dir?
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
