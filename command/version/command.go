// Package version implements the version command for the command line tool.
package version

import (
	"fmt"
	"runtime"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

type Config struct {
	Description string
	Name        string
	Source      string
}

type Command struct {
	cobraCommand *cobra.Command

	description string
	name        string
	source      string
}

func New(config Config) (*Command, error) {
	if config.Description == "" {
		return nil, errors.Errorf("%T.Description must not be empty", config)
	}
	if config.Name == "" {
		return nil, errors.Errorf("%T.Name must not be empty", config)
	}
	if config.Source == "" {
		return nil, errors.Errorf("%T.Source must not be empty", config)
	}

	c := &Command{
		cobraCommand: nil,

		description: config.Description,
		name:        config.Name,
		source:      config.Source,
	}

	c.cobraCommand = &cobra.Command{
		Use:   "version",
		Short: "Show version information of the command line tool.",
		Long:  "Show version information of the command line tool.",
		Run:   c.Execute,
	}

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	fmt.Printf("Description:    %s\n", c.description)
	fmt.Printf("Go Version:     %s\n", runtime.Version())
	fmt.Printf("Name:           %s\n", c.name)
	fmt.Printf("OS / Arch:      %s / %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Source:         %s\n", c.source)
}
