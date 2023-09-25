package cmd

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/harrisonde/adel"
	"github.com/rodaine/table"
)

type CommandsHelper map[int]adel.Command

//go:embed cli/*.go
var efs embed.FS

func GetHelp() string {

	var commands = CommandsHelper{}

	adelCommands, err := LoadDefaultCommands()
	if err != nil {
		color.Yellow(fmt.Sprintf("error loading adel help command, %s", err))
	}

	p := 1
	for _, command := range adelCommands {
		commands[p] = command
		p++
	}

	appCommands, err := LoadCommands()
	if err != nil {
		color.Yellow(fmt.Sprintf("error loading help, you may have a malformed command in your command dir; %s", err))
	}

	p++

	for _, c := range appCommands {
		commands[p] = c
		p++
	}

	// Sort
	keys := make([]int, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return commands[keys[i]].Name < commands[keys[j]].Name
	})

	// table
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Command", "Description")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, k := range keys {
		tbl.AddRow(commands[k].Name, commands[k].Help)
	}

	color.Yellow("Available commands:\n\n")
	tbl.Print()

	return ""
}

func LoadDefaultCommands() (CommandsHelper, error) {

	var commandsHelpers = CommandsHelper{}
	pointer := 1
	fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		body, err := efs.ReadFile(path)
		if err != nil {
			return err
		}

		commandBody := string(body)

		cmd, _ := ParseCommand(commandBody)

		if *cmd != (adel.Command{}) {
			commandsHelpers[pointer] = *cmd
		}

		pointer++
		return nil
	})
	return commandsHelpers, nil
}

func LoadCommands() (CommandsHelper, error) {

	var commandsHelpers = CommandsHelper{}

	path, err := os.Getwd()
	if err != nil {
		color.Red(fmt.Sprintf("%s", err))
	}
	path = path + "/cmd"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		color.Red(fmt.Sprintf("%s", err))
	}

	for index, file := range files {
		if file.IsDir() {
			color.Yellow("commands may only contain files, folders are skipped")
		}

		body, err := ioutil.ReadFile(path + "/" + file.Name())
		if err != nil {
			return commandsHelpers, err
		}

		commandBody := string(body)

		cmd, _ := ParseCommand(commandBody)

		if *cmd != (adel.Command{}) {
			commandsHelpers[index] = *cmd
		}

	}
	return commandsHelpers, nil
}

func ParseCommand(commandBody string) (*adel.Command, error) {

	cmd := &adel.Command{}
	r := regexp.MustCompile("&adel.Command{(.*\n.*\n.*\n.*)")

	bodyMap := r.FindStringSubmatch(commandBody)
	if len(bodyMap) != 2 {
		return cmd, errors.New("malformed command; please use make command to create a new command template")
	}
	st := bodyMap[1]
	st = strings.Replace(st, "\t", "", -1)
	st = strings.Replace(st, "\n", "", -1)

	var name string
	var helper string

	res := strings.Split(st, ",")
	for _, c := range res {
		m := strings.Split(c, ":")

		if len(m) > 1 {
			key := strings.ToLower(m[0])
			value := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(m[1], "")
			value = strings.TrimSpace(value)
			switch key {
			case "name":
				name = value
			case "help":
				helper = value
			}
		}
	}

	if name != "" {
		cmd.Name = name
	}
	if helper != "" {
		cmd.Help = helper
	}

	return cmd, nil
}
