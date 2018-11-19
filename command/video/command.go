package video

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/marcelmue/angory/command/video/annotate"
	"github.com/marcelmue/angory/command/video/pull"
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
		Use:   "video",
		Short: "Do all the things with video metadata, e.g. pull them from youtube.",
		Long:  "Do all the things with video metadata, e.g. pull them from youtube.",
		Run:   c.Execute,
	}

	var annotateCommand *annotate.Command
	{
		annotateCommand, err = annotate.New()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var pullCommand *pull.Command
	{
		pullCommand, err = pull.New()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	c.cobraCommand.AddCommand(annotateCommand.CobraCommand())
	c.cobraCommand.AddCommand(pullCommand.CobraCommand())

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
}
