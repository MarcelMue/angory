package metadata

import (
	"github.com/juju/errors"
	"github.com/marcelmue/angory/pkg/annotation"
	"github.com/marcelmue/angory/pkg/entities"
	"github.com/marcelmue/angory/pkg/game"
	"github.com/marcelmue/angory/pkg/talent"
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

func (s *Service) AnnotateVideos() ([]*entities.Video, error) {
	games, err := game.FromPath(s.gamesPath)
	if err != nil {
		return []*entities.Video{}, errors.Trace(err)
	}
	gamesMap := gamesMap(games)

	talents, err := talent.FromPath(s.talentsPath)
	if err != nil {
		return []*entities.Video{}, errors.Trace(err)
	}
	talentsMap := talentsMap(talents)

	annotations, err := annotation.FromPath(s.videoAnnotationsPath)
	if err != nil {
		return []*entities.Video{}, errors.Trace(err)
	}
	annotationsMap := annotationsMap(annotations, gamesMap, talentsMap)

	videos, err := entities.FromYoutubeVideosPath(s.youtubeVideosPath)
	if err != nil {
		return []*entities.Video{}, errors.Trace(err)
	}

	for _, video := range videos {
		if val, ok := annotationsMap[video.ID]; ok {
			video.Annotation = &val
		}
	}
	return videos, nil
}

func annotationsMap(annotations []annotation.Annotation, games map[string]entities.Game, talents map[string]entities.Talent) map[string]entities.Annotation {
	m := make(map[string]entities.Annotation)
	for _, annotation := range annotations {
		var annotatedTalents []*entities.Talent
		for _, talentID := range annotation.TalentIDs {
			if val, ok := talents[talentID]; ok {
				annotatedTalents = append(annotatedTalents, &val)
			}
		}

		var annotatedGame *entities.Game
		if val, ok := games[annotation.GameID]; ok {
			annotatedGame = &val
		}

		m[annotation.VideoID] = entities.Annotation{
			Game:    annotatedGame,
			Talents: annotatedTalents,
		}
	}
	return m
}

func gamesMap(games []game.Game) map[string]entities.Game {
	m := make(map[string]entities.Game)
	for _, game := range games {
		m[game.ID] = entities.Game{
			Name:  game.Name,
			Steam: game.Steam,
			Itch:  game.Itch,
		}
	}
	return m
}

func talentsMap(talents []talent.Talent) map[string]entities.Talent {
	m := make(map[string]entities.Talent)
	for _, talent := range talents {
		m[talent.ID] = entities.Talent{
			Channel:  talent.Channel,
			Name:     talent.Name,
			Nickname: talent.Nickname,
			Short:    talent.Short,
			Twitter:  talent.Twitter,
		}
	}
	return m
}
