package main

import (
	"fmt"

	"github.com/cidekar/adele-framework"
	"github.com/fatih/color"
)

var VersionCommand = &Command{
	Name:        "version",
	Help:        "Print the current version",
	Description: "Print the current framework version to the screen",
	Usage:       "adele version [options]",
	Examples:    []string{"adele version", "adele version --verbose"},
	Options:     map[string]string{},
}

type CommandVersion interface {
	Handle() error
	Get() string
}

type Version struct {
	CommandVersion
}

func NewVersion() *Version {
	return &Version{}
}

func (c *Version) Handle() error {
	version := color.GreenString(adele.Version)
	fmt.Println("Adele Framework", version)
	return nil
}

func init() {
	if err := Registry.Register(VersionCommand); err != nil {
		panic(fmt.Sprintf("Failed to register version command: %v", err))
	}
}
