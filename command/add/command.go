package add

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/marcelmue/angory/command/add/talent"
)

type Command struct {
	cobraCommand *cobra.Command
}

// New creates a new list command.
func New() (*Command, error) {
	var err error

	c := &Command{
		cobraCommand: nil,
	}

	c.cobraCommand = &cobra.Command{
		Use:   "add",
		Short: "Add new metadata, e.g. a new talent.",
		Long:  "Add new metadata, e.g. a new talent.",
		Run:   c.Execute,
	}

	var talentCommand *talent.Command
	{
		talentCommand, err = talent.New()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	c.cobraCommand.AddCommand(talentCommand.CobraCommand())

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
}
