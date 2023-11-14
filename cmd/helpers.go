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
	"github.com/harrisonde/adele-framework"
	"github.com/rodaine/table"
)

type CommandsHelper map[int]adele.Command

//go:embed cli/*.go
var efs embed.FS

// use in myapp cmd with exiting a command
func Exit(message string) string {
	return message
}

// use in myapp cmd with exiting a command with error
func ExitError(err error) {
	color.Red("error: %v", err)
}

func GetHelp(args ...string) string {

	var commands = CommandsHelper{}

	adeleCommands, err := LoadDefaultCommands()
	if err != nil {
		color.Yellow(fmt.Sprintf("error loading adele help command, %s", err))
	}

	p := 1
	for _, command := range adeleCommands {
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

	// Print a single command
	if len(args) > 0 {
		printCommandTable(commands, args[0])
	} else {
		printHelpTable(commands)
	}
	return ""
}

func printOptionsTable(cmdOptions map[string]string) {

	options := map[string]string{
		"-h, --help":    "display help for the given command.",
		"-v, --verbose": "increase the verbosity of messages",
	}

	headerFmt := color.New(color.FgYellow).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tblOpts := table.New("Options:", "")
	tblOpts.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for k, v := range cmdOptions {
		tblOpts.AddRow(fmt.Sprintf("  %s", k), strings.Trim(v, " "))
	}

	for k, v := range options {
		tblOpts.AddRow(fmt.Sprintf("  %s", k), strings.Trim(v, " "))
	}
	tblOpts.Print()
}

func printCommandTable(commands CommandsHelper, arg1 string) {
	commandFound := false

	for _, c := range commands {
		if c.Name == arg1 {
			commandFound = true
			color.Yellow("Description:")
			fmt.Printf(" " + c.Description + "\n\n")
			color.Yellow("Usage:")
			fmt.Printf(" " + c.Usage + "\n\n")
			printOptionsTable(c.Options)
		}
	}

	if !commandFound {
		printHelpTable(commands)
	}
}

func printHelpTable(commands CommandsHelper) {
	// Sort and print commands
	keys := make([]int, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return commands[keys[i]].Name < commands[keys[j]].Name
	})

	// table
	headerFmt := color.New(color.FgYellow).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()

	tblCmd := table.New("Available commands:", "")
	tblCmd.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, k := range keys {
		tblCmd.AddRow("  "+commands[k].Name, commands[k].Help)
	}

	tblOpts := table.New("Options:", "")
	tblOpts.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tblOpts.AddRow("  -h, --help", "display help for the given command.")

	fmt.Printf("Adele framework\n")
	color.Yellow("\nUsage:")
	fmt.Printf("  command [options] [arguments]\n\n")
	tblOpts.Print()
	fmt.Printf("\n")
	tblCmd.Print()
}

func LoadDefaultCommands() (CommandsHelper, error) {
	var commandsHelpers = CommandsHelper{}
	excludeFiles := []string{
		"cli/copy-files.go",
		"cli/helpers.go",
		"cli/main.go",
		"cli/rpc.go",
	}
	pointer := 1
	fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		for _, e := range excludeFiles {
			if e == path {
				return nil
			}
		}

		body, err := efs.ReadFile(path)
		if err != nil {
			return err
		}

		commandBody := string(body)

		cmd, err := ParseCommand(commandBody)
		if err == nil {
			commandsHelpers[pointer] = *cmd
		} else {
			color.Red(fmt.Sprintf("%s %s", path, err))
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
		if err == nil {
			commandsHelpers[index] = *cmd
		} else {
			color.Red(fmt.Sprintf("error loading %s, %s", path, err))
		}

	}
	return commandsHelpers, nil
}

func ParseCommand(commandBody string) (*adele.Command, error) {

	cmd := &adele.Command{}
	r := regexp.MustCompile("&adele.Command{(?s)(.*})")

	cmdStruct := r.FindStringSubmatch(commandBody)

	if len(cmdStruct) != 2 {
		return cmd, errors.New("malformed command; please use make command to create a new command template")
	}
	st := cmdStruct[1]
	st = strings.Replace(st, "\t", "", -1)
	st = strings.Replace(st, "\n", "", -1)

	res := strings.Split(st, "\",")
	for _, c := range res {
		m := strings.Split(c, ":")

		if len(m) > 1 {
			key := strings.ToLower(m[0])
			value := strings.TrimSpace(m[1])

			switch key {
			case "name":
				cmd.Name = strings.ReplaceAll(value, `"`, "")
			case "options":
				cmd.Options = map[string]string{}
				r = regexp.MustCompile("Options:(.*?)}")
				opts := r.FindStringSubmatch(st)

				if len(opts) == 2 {
					os := strings.Split(opts[1], "\",")

					for _, v := range os {
						v = strings.Replace(v, "map[string]string{", "", 1)
						o := strings.Split(v, ":")
						if len(o) == 2 {
							flags := strings.Replace(o[0], "map[string]string{", "", 1)
							flags = strings.ReplaceAll(flags, `"`, "")
							flags = strings.TrimSpace(flags)
							description := strings.ReplaceAll(o[1], `"`, "")
							description = strings.ReplaceAll(description, `"`, "")
							description = strings.TrimSpace(description)
							cmd.Options[flags] = description
						}
					}
				}

			case "usage":
				cmd.Usage = strings.ReplaceAll(value, `"`, "")

			case "description":
				cmd.Description = strings.ReplaceAll(value, `"`, "")
			case "help":
				cmd.Help = strings.ReplaceAll(value, `"`, "")
			}
		}
	}

	if cmd == (&adele.Command{}) {
		return cmd, errors.New("empty command; please use make command to create a new command template")
	}

	return cmd, nil
}
