package metadata

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
	"github.com/marcelmue/angory/pkg/annotation"
	"github.com/marcelmue/angory/pkg/talent"
	"github.com/marcelmue/angory/service/metadata/video"
	youtube "google.golang.org/api/youtube/v3"
)

type Config struct {
	TalentsPath          string
	VideoAnnotationsPath string
	YoutubeVideosPath    string
}
type Service struct {
	talentsPath          string
	videoAnnotationsPath string
	youtubeVideosPath    string
}

func New(config Config) (*Service, error) {
	if config.TalentsPath == "" {
		return nil, errors.Errorf("%T.TalentsPath must not be empty", config)
	}
	if config.VideoAnnotationsPath == "" {
		return nil, errors.Errorf("%T.AnnotationsPath must not be empty", config)
	}
	if config.YoutubeVideosPath == "" {
		return nil, errors.Errorf("%T.YoutubeVideosPath must not be empty", config)
	}

	s := &Service{
		talentsPath:          config.TalentsPath,
		videoAnnotationsPath: config.VideoAnnotationsPath,
		youtubeVideosPath:    config.YoutubeVideosPath,
	}

	return s, nil
}

func (s *Service) AnnotateVideos() ([]*video.Video, error) {
	annotations, err := annotation.FromPath(s.videoAnnotationsPath)
	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}
	talents, err := talent.FromPath(s.talentsPath)
	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}
	talentsMap := talentsMap(talents)
	annotationsMap := annotationsMap(annotations, talentsMap)
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

func annotationsMap(annotations []annotation.Annotation, talents map[string]video.Talent) map[string]video.Annotation {
	m := make(map[string]video.Annotation)
	for _, annotation := range annotations {
		annotatedTalents := []video.Talent{}
		for _, talentID := range annotation.TalentIDs {
			if val, ok := talents[talentID]; ok {
				annotatedTalents = append(annotatedTalents, val)
			}
		}
		m[annotation.VideoID] = video.Annotation{
			Game:    annotation.Game,
			Talents: annotatedTalents,
		}
	}
	return m
}

func talentsMap(talents []talent.Talent) map[string]video.Talent {
	m := make(map[string]video.Talent)
	for _, talent := range talents {
		m[talent.ID] = video.Talent{
			Channel:  talent.Channel,
			Name:     talent.Name,
			Nickname: talent.Nickname,
			Short:    talent.Short,
			Twitter:  talent.Twitter,
		}
	}
	return m
}
