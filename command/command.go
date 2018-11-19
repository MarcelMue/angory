package command

import (
	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/marcelmue/angory/command/add"
	"github.com/marcelmue/angory/command/version"
	"github.com/marcelmue/angory/command/video"
)

type Config struct {
	Description string
	Name        string
	Source      string
}

type Command struct {
	cobraCommand *cobra.Command
}

func New(config Config) (*Command, error) {
	var err error

	c := &Command{
		cobraCommand: nil,
	}

	c.cobraCommand = &cobra.Command{
		Use:   config.Name,
		Short: config.Description,
		Long:  config.Description,
		Run:   c.Execute,
	}

	var addCommand *add.Command
	{
		addCommand, err = add.New()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var videoCommand *video.Command
	{
		videoCommand, err = video.New()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	var versionCommand *version.Command
	{
		c := version.Config{
			Description: config.Description,
			Name:        config.Name,
			Source:      config.Source,
		}

		versionCommand, err = version.New(c)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	c.cobraCommand.AddCommand(addCommand.CobraCommand())
	c.cobraCommand.AddCommand(videoCommand.CobraCommand())
	c.cobraCommand.AddCommand(versionCommand.CobraCommand())

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	cmd.HelpFunc()(cmd, nil)
}
