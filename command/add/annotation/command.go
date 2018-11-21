package annotation

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/juju/errors"
	"github.com/marcelmue/angory/pkg/annotation"
	"github.com/marcelmue/angory/pkg/entities"
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
		Use:   "annotation",
		Short: "Add a new annotation to the set of metadata.",
		Long:  "Add a new talent to the set of metadata. The output is added to the JSON file.",
		Run:   c.Execute,
	}

	c.cobraCommand.PersistentFlags().BoolVarP(&flags.Missing, "missing", "m", false, "Add an annotation to a video that doesn't have one yet.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.AnnotationsPath, "path", "p", "data/metadata/video_annotations.json", "Path to read and write the talents.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.VideoAnnotationsPath, "video-annotations", "", "data/metadata/video_annotations.json", "Path to the annotations for youtube videos.")
	c.cobraCommand.PersistentFlags().StringVarP(&flags.YoutubeVideosPath, "yt-videos", "", "data/output/youtube_videos.json", "Path to the pulled youtube video metadata.")

	return c, nil
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	newAnnotation := annotation.Annotation{}
	var err error
	reader := bufio.NewReader(os.Stdin)
	{
		if flags.Missing {
			newAnnotation.VideoID, err = findMissing()
			if err != nil {
				log.Fatalf("Error: %v", errors.ErrorStack(err))
				os.Exit(1)
			}

			if newAnnotation.VideoID == "" {
				log.Fatalf("All videos are already annotated")
				os.Exit(1)
			} else {
				log.Printf("You are annotating video (%s) link: https://www.youtube.com/watch?v=%s \n", newAnnotation.VideoID, newAnnotation.VideoID)
			}

		} else {
			fmt.Println("Enter the annotations VideoID:")
			newAnnotation.VideoID, err = reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error: %v", errors.ErrorStack(err))
				os.Exit(1)
			}
			newAnnotation.VideoID = strings.Replace(newAnnotation.VideoID, "\n", "", -1)
			if newAnnotation.VideoID == "" {
				log.Fatal("VideoID can not be empty")
				os.Exit(1)
			}
		}
	}

	{
		fmt.Println("Enter the annotations GameID (one word: no capitals, multiple words: every starting letter):")
		newAnnotation.GameID, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		newAnnotation.GameID = strings.Replace(newAnnotation.GameID, "\n", "", -1)
	}

	{
		fmt.Println("Enter the annotations talents (separated by commas):")
		talents, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error: %v", errors.ErrorStack(err))
			os.Exit(1)
		}
		talents = strings.Replace(talents, "\n", "", -1)
		newAnnotation.TalentIDs = strings.Split(talents, ",")
	}

	err = annotation.ToPath(flags.AnnotationsPath, newAnnotation)
	if err != nil {
		log.Fatalf("Error: %v", errors.ErrorStack(err))
		os.Exit(1)
	}

	log.Printf("Added annotation (%s) to the file %s \n", newAnnotation.VideoID, flags.AnnotationsPath)
}

func annotationIDs(annotations []annotation.Annotation) map[string]string {
	m := make(map[string]string)
	for _, annotation := range annotations {
		m[annotation.VideoID] = annotation.VideoID
	}
	return m
}

func findMissing() (string, error) {
	annotations, err := annotation.FromPath(flags.VideoAnnotationsPath)
	if err != nil {
		return "", errors.Trace(err)
	}
	annotationIDs := annotationIDs(annotations)
	videos, err := entities.FromYoutubeVideosPath(flags.YoutubeVideosPath)
	if err != nil {
		return "", errors.Trace(err)
	}
	for _, video := range videos {
		_, ok := annotationIDs[video.ID]
		if !ok {
			return video.ID, nil
		}
	}

	return "", nil
}
