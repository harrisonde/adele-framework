package main

type CliInterface interface {
	Handle(string, string, string, string) error
	Help(string string) error
	Validate() (string, string, string, string, error)
}

type Cli struct {
	CliInterface
}

type Command struct {
	Name        string            `cmd:"name" help:"The command name" required:"true"`
	Description string            `cmd:"description" help:"Brief description of the command"`
	Usage       string            `cmd:"usage" help:"Usage pattern for the command"`
	Help        string            `cmd:"help" help:"Short help text"`
	Examples    []string          `cmd:"examples" help:"Usage examples"`
	Options     map[string]string `cmd:"options" help:"Available command options"`
}

type CommandsHelper map[int]Command
