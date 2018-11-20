package annotate

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/marcelmue/angory/service/metadata"
	"github.com/marcelmue/angory/service/metadata/video"
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
		Use:   "annotate",
		Short: "Annotate youtube videos with additional metadata and write it to a file.",
		Long:  "Annotate youtube videos with additional metadata and write it to a file. The output is a JSON.",
		Run:   c.Execute,
	}

	c.cobraCommand.PersistentFlags().StringVarP(&flags.AnnotatedVideosPath, "annotated-videos", "p", "data/output/annotated_videos.json", "Path to write the output file to.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.GamesPath, "games", "", "data/metadata/games.json", "Path to the games for video annotations.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.TalentsPath, "talents", "", "data/metadata/talents.json", "Path to the talents for video annotations.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.VideoAnnotationsPath, "video-annotations", "", "data/metadata/video_annotations.json", "Path to the annotations for youtube videos.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.YoutubeVideosPath, "yt-videos", "", "data/output/youtube_videos.json", "Path to the pulled youtube video metadata.")

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	var annotatedVideos []*video.Video
	{
		config := metadata.Config{
			GamesPath:            flags.GamesPath,
			TalentsPath:          flags.TalentsPath,
			VideoAnnotationsPath: flags.VideoAnnotationsPath,
			YoutubeVideosPath:    flags.YoutubeVideosPath,
		}
		s, err := metadata.New(config)
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		annotatedVideos, err = s.AnnotateVideos()
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
	}

	videosJSON, err := json.MarshalIndent(annotatedVideos, "", "  ")
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}
	err = os.MkdirAll(path.Dir(flags.AnnotatedVideosPath), os.ModePerm)
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}

	err = ioutil.WriteFile(flags.AnnotatedVideosPath, videosJSON, 0644)
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}

	log.Printf("Annotated the complete list of all videos from %s with annotations from %s to file %s \n", flags.YoutubeVideosPath, flags.VideoAnnotationsPath, flags.AnnotatedVideosPath)
}
