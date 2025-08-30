package main

import "fmt"

type CommandRegistry struct {
	commands map[string]*Command
}

var Registry = NewCommandRegistry()

// NewCommandRegistry creates and returns a new CommandRegistry with an initialized command map.
// Use this constructor to create new registry instances for testing or isolated command sets.
//
// Example:
//
//	registry := NewCommandRegistry()
//	registry.Register(myCommand)
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*Command),
	}
}

// Register adds a command to the registry after validation. Returns an error if the command
// fails validation checks (missing required fields, invalid format, etc.). Commands with
// duplicate names will overwrite existing entries.
//
// Example:
//
//	err := registry.Register(&Command{Name: "version", Description: "Show version"})
//	if err != nil {
//	    log.Fatal("Failed to register command:", err)
//	}
func (cr *CommandRegistry) Register(cmd *Command) error {
	if err := cr.validateCommand(cmd); err != nil {
		return fmt.Errorf("invalid command %s: %w", cmd.Name, err)
	}
	cr.commands[cmd.Name] = cmd
	return nil
}

// GetCommand retrieves a command by name from the registry. Returns the command pointer
// and a boolean indicating whether the command was found. The boolean should always be
// checked before using the returned command pointer.
//
// Example:
//
//	if cmd, exists := registry.GetCommand("version"); exists {
//	    fmt.Println(cmd.Description)
//	}
func (cr *CommandRegistry) GetCommand(name string) (*Command, bool) {
	cmd, exists := cr.commands[name]
	return cmd, exists
}

// GetAllCommands returns a map of all registered commands keyed by command name.
// The returned map is a direct reference to the internal storage - modifications
// will affect the registry. Use with caution in concurrent environments.
//
// Example:
//
//	for name, cmd := range registry.GetAllCommands() {
//	    fmt.Printf("%s: %s\n", name, cmd.Description)
//	}
func (cr *CommandRegistry) GetAllCommands() map[string]*Command {
	return cr.commands
}
