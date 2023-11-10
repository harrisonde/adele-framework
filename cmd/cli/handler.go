package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/harrisonde/adele"
	"github.com/iancoleman/strcase"
)

var MakeHandlerCommand = &adele.Command{
	Name:        "make handler",
	Help:        "create a new handler",
	Description: "use the make handler command to create a new handler in the handlers directory",
	Usage:       "make handler <name>",
	Options:     map[string]string{},
}

func doMakeHandler(arg3 string) error {
	if arg3 == "" {
		return errors.New("you must provide a name for the handler")
	}

	fileName := ade.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
	if fileExists(fileName) {
		return errors.New(fileName + " already exits!")
	}

	data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
	if err != nil {
		return err
	}

	// Read the template and find/replace with the provided handler name
	handler := string(data)
	handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

	// Write the file
	err = ioutil.WriteFile(fileName, []byte(handler), 0644)
	if err != nil {
		return err
	}

	return nil
}
