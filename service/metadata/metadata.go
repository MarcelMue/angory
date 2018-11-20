package metadata

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
	"github.com/marcelmue/angory/pkg/annotation"
	"github.com/marcelmue/angory/pkg/game"
	"github.com/marcelmue/angory/pkg/talent"
	"github.com/marcelmue/angory/service/metadata/video"
	youtube "google.golang.org/api/youtube/v3"
)

type Config struct {
	GamesPath            string
	TalentsPath          string
	VideoAnnotationsPath string
	YoutubeVideosPath    string
}
type Service struct {
	gamesPath            string
	talentsPath          string
	videoAnnotationsPath string
	youtubeVideosPath    string
}

func New(config Config) (*Service, error) {
	if config.GamesPath == "" {
		return nil, errors.Errorf("%T.GamesPath must not be empty", config)
	}
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
		gamesPath:            config.GamesPath,
		talentsPath:          config.TalentsPath,
		videoAnnotationsPath: config.VideoAnnotationsPath,
		youtubeVideosPath:    config.YoutubeVideosPath,
	}

	return s, nil
}

func (s *Service) AnnotateVideos() ([]*video.Video, error) {
	games, err := game.FromPath(s.gamesPath)
	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}
	gamesMap := gamesMap(games)

	talents, err := talent.FromPath(s.talentsPath)
	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}
	talentsMap := talentsMap(talents)

	annotations, err := annotation.FromPath(s.videoAnnotationsPath)
	if err != nil {
		return []*video.Video{}, errors.Trace(err)
	}
	annotationsMap := annotationsMap(annotations, gamesMap, talentsMap)

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

func annotationsMap(annotations []annotation.Annotation, games map[string]video.Game, talents map[string]video.Talent) map[string]video.Annotation {
	m := make(map[string]video.Annotation)
	for _, annotation := range annotations {
		var annotatedTalents []*video.Talent
		for _, talentID := range annotation.TalentIDs {
			if val, ok := talents[talentID]; ok {
				annotatedTalents = append(annotatedTalents, &val)
			}
		}

		var annotatedGame *video.Game
		if val, ok := games[annotation.GameID]; ok {
			annotatedGame = &val
		}

		m[annotation.VideoID] = video.Annotation{
			Game:    annotatedGame,
			Talents: annotatedTalents,
		}
	}
	return m
}

func gamesMap(games []game.Game) map[string]video.Game {
	m := make(map[string]video.Game)
	for _, game := range games {
		m[game.ID] = video.Game{
			Name:  game.Name,
			Steam: game.Steam,
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
