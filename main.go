package main

import (
	"github.com/marcelmue/angory/command"
)

const angoryChannel = "UC5rUMdCFWPXYs9e8PBLzq5g"

var (
	description string = "Tool to pull and enhance metadata for youtube channels."
	name        string = "angory"
	source      string = "https://github.com/marcelmue/angory"
)

func main() {
	var err error

	var newCommand *command.Command
	{
		c := command.Config{
			Description: description,
			Name:        name,
			Source:      source,
		}

		newCommand, err = command.New(c)
		if err != nil {
			panic(err)
		}
	}

	newCommand.CobraCommand().Execute()
}
