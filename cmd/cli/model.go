package main

import (
	"errors"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/harrisonde/adele"
	"github.com/iancoleman/strcase"
)

var MakeModelCommand = &adele.Command{
	Name:        "make model",
	Help:        "create a new model",
	Description: "use the make model command to create a new model in the models directory",
	Usage:       "make model <name>",
	Options:     map[string]string{},
}

func doMakeModel(arg3 string) error {
	if arg3 == "" {
		return errors.New("You must provide a name for the model!")
	}

	data, err := templateFS.ReadFile("templates/data/model.go.txt")
	if err != nil {
		return err
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
		return errors.New(fileName + " already exits!")
	}

	model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)
	err = copyDataToFile([]byte(model), fileName)
	if err != nil {
		return err
	}
	return nil
}
