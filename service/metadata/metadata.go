package metadata

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
	annotation "github.com/marcelmue/angory/service/metadata/Annotation"
	"github.com/marcelmue/angory/service/metadata/video"
	youtube "google.golang.org/api/youtube/v3"
)

type Config struct {
	VideoAnnotationsPath string
	YoutubeVideosPath    string
}
type Service struct {
	videoAnnotationsPath string
	youtubeVideosPath    string
}

func New(config Config) (*Service, error) {
	if config.VideoAnnotationsPath == "" {
		return nil, errors.Errorf("%T.AnnotationsPath must not be empty", config)
	}
	if config.YoutubeVideosPath == "" {
		return nil, errors.Errorf("%T.YoutubeVideosPath must not be empty", config)
	}

	s := &Service{
		videoAnnotationsPath: config.VideoAnnotationsPath,
		youtubeVideosPath:    config.YoutubeVideosPath,
	}

	return s, nil
}

func (s *Service) AnnotateVideos() ([]*video.Video, error) {
	annotations, err := s.fromAnnotationsPath()
	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}
	annotationsMap := annotationsMap(annotations)
	videos, err := s.fromYoutubeVideosPath()

	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}
	for _, video := range videos {
		if val, ok := annotationsMap[video.ID]; ok {
			video.Annotation = &val
		}
	}
	return videos, nil
}

func (s *Service) fromAnnotationsPath() ([]annotation.Annotation, error) {
	jsonFile, err := os.Open(s.videoAnnotationsPath)
	if err != nil {
		return []annotation.Annotation{}, errors.Trace(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []annotation.Annotation{}, errors.Trace(err)
	}

	var result []annotation.Annotation
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

func (s *Service) fromYoutubeVideosPath() ([]*video.Video, error) {
	jsonFile, err := os.Open(s.youtubeVideosPath)
	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}

	var youtubeVideos []*youtube.Video
	json.Unmarshal([]byte(byteValue), &youtubeVideos)

	var result []*video.Video
	for _, ytVideo := range youtubeVideos {
		vid := &video.Video{
			ChannelID:   ytVideo.Snippet.ChannelId,
			Description: ytVideo.Snippet.Description,
			ID:          ytVideo.Id,
			Kind:        ytVideo.Kind,
			PublishedAt: ytVideo.Snippet.PublishedAt,
			Title:       ytVideo.Snippet.Title,
			Statistics: video.Statistics{
				CommentCount: int(ytVideo.Statistics.CommentCount),
				DislikeCount: int(ytVideo.Statistics.DislikeCount),
				LikeCount:    int(ytVideo.Statistics.LikeCount),
				ViewCount:    int(ytVideo.Statistics.ViewCount),
			},
			YoutubeTags: ytVideo.Snippet.Tags,
		}
		result = append(result, vid)
	}

	return result, nil
}

func annotationsMap(annotations []annotation.Annotation) map[string]video.Annotation {
	m := make(map[string]video.Annotation)
	for _, annotation := range annotations {
		m[annotation.VideoID] = video.Annotation{
			Game:   annotation.Game,
			Talent: annotation.Talent,
		}
	}
	return m
}
