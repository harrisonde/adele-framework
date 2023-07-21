package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3, arg4 string) error {

	switch arg2 {
	case "key":
		rnd := ade.RandomString(32)
		color.Yellow("Encryption key: %s", rnd)

	case "migration":

		checkForDb()

		if arg3 == "" {
			exitGracefully(errors.New("You must provide a name for the migration!"))
		}

		// Default to fizz migrations
		migrationType := "fizz"
		var up, down string

		// What type of migration?
		if arg4 == "fizz" || arg4 == "" {
			upBytes, err := templateFS.ReadFile("templates/migrations/migration_up.fizz")
			if err != nil {
				exitGracefully(err)
			}

			downBytes, err := templateFS.ReadFile("templates/migrations/migration_down.fizz")
			if err != nil {
				exitGracefully(err)
			}

			up = string(upBytes)
			down = string(downBytes)
		} else {
			migrationType = "sql"
		}

		err := ade.CreatePopMigration([]byte(up), []byte(down), arg3, migrationType)
		if err != nil {
			exitGracefully(err)
		}

	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}

	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("You must provide a name for the handler!"))
		}

		fileName := ade.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exits!"))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		// Read the template and find/replace with the provided handler name
		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		// Write the file
		err = ioutil.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}

		return nil

	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("You must provide a name for the model!"))
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		// Read the model and find/replace with the provided handler name
		model := string(data)
		plur := pluralize.NewClient()

		var modelName = arg3
		var tableName = arg3

		// Get the name provided into the format we need them to do
		// i.e., table name is plural model is singular
		if plur.IsPlural(arg3) {
			modelName = plur.Singular(arg3)
			tableName = strings.ToLower(arg3)
		} else {
			tableName = strings.ToLower(plur.Plural(arg3))
		}

		// Build up the file and replace with model and table name
		fileName := ade.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exits!"))
		}
		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			exitGracefully(err)
		}
	case "mail":
		if arg3 == "" {
			exitGracefully(errors.New("You must provide a name for the mail template!"))
		}

		htmlMail := ade.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainMail := ade.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

		err := copyFileFromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFileFromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
		if err != nil {
			exitGracefully(err)
		}

	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}
