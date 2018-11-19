package pull

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"

	"github.com/marcelmue/angory/service/playlist"
	"github.com/marcelmue/angory/service/video"
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
		Use:   "pull",
		Short: "Pull video metadata from the youtube API for a channel and write it to a file.",
		Long:  "Pull video metadata from the youtube API for a channel and write it to a file. The output is a JSON.",
		Run:   c.Execute,
	}

	c.cobraCommand.PersistentFlags().StringVarP(&flags.APIKey, "key", "k", "", "API key for the youtube v3 API.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.Path, "path", "p", "data/youtube_videos.json", "Path to write the output file to.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.YoutubeChannel, "channel", "c", "UC5rUMdCFWPXYs9e8PBLzq5g", "ID of the youtube channel to pull videos from.")

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	err := flags.Validate()
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}

	var youtubeClient *youtube.Service
	{
		client := &http.Client{
			Transport: &transport.APIKey{Key: flags.APIKey},
		}
		youtubeClient, err = youtube.New(client)
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
	}

	var playlistID string
	{
		config := playlist.Config{
			YoutubeClient: youtubeClient,
		}
		s, err := playlist.New(config)
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		playlistID, err = s.GetUploadID(flags.YoutubeChannel)
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
	}

	var videos []*youtube.Video
	{
		config := video.Config{
			YoutubeClient: youtubeClient,
		}
		s, err := video.New(config)
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		videos, err = s.FromPlaylist(playlistID)
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
	}

	videosJSON, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}
	err = os.MkdirAll(path.Dir(flags.Path), os.ModePerm)
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}

	err = ioutil.WriteFile(flags.Path, videosJSON, 0644)
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}

	log.Printf("Pulled the complete list of all videos for %s to file %s \n", flags.YoutubeChannel, flags.Path)

}
