package talent

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/marcelmue/angory/pkg/talent"
	"github.com/spf13/cobra"
)

var (
	flags = &Flags{}
)

type Command struct {
	cobraCommand *cobra.Command
}

func New() (*Command, error) {

	c := &Command{
		cobraCommand: nil,
	}

	c.cobraCommand = &cobra.Command{
		Use:   "talent",
		Short: "Add a new talent to the set of metadata.",
		Long:  "Add a new talent to the set of metadata. The output is added to the JSON file.",
		Run:   c.Execute,
	}

	c.cobraCommand.PersistentFlags().StringVarP(&flags.TalentsPath, "path", "p", "data/metadata/talents.json", "Path to read and write the talents.")

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	newTalent := talent.Talent{}
	var err error
	reader := bufio.NewReader(os.Stdin)
	{
		fmt.Println("Enter the talents ID:")
		newTalent.ID, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		newTalent.ID = strings.Replace(newTalent.ID, "\n", "", -1)
		if newTalent.ID == "" {
			log.Fatal("ID can not be empty")
			os.Exit(1)
		}
	}
	{
		fmt.Println("Enter the talents Channel:")
		newTalent.Channel, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		newTalent.Channel = strings.Replace(newTalent.Channel, "\n", "", -1)
	}
	{
		fmt.Println("Enter the talents Name:")
		newTalent.Name, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		newTalent.Name = strings.Replace(newTalent.Name, "\n", "", -1)
	}

	{
		fmt.Println("Enter the talents Nickname:")
		newTalent.Nickname, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		newTalent.Nickname = strings.Replace(newTalent.Nickname, "\n", "", -1)
	}

	{
		fmt.Println("Enter the talents Short hand:")
		newTalent.Short, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		newTalent.Short = strings.Replace(newTalent.Short, "\n", "", -1)
	}

	{
		fmt.Println("Enter the talents Twitter:")
		newTalent.Twitter, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		newTalent.Twitter = strings.Replace(newTalent.Twitter, "\n", "", -1)
	}

	err = talent.ToPath(flags.TalentsPath, newTalent)
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}

	log.Printf("Added talent (%s) to the file %s \n", newTalent.ID, flags.TalentsPath)
}
