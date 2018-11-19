package playlist

import (
	"github.com/juju/errors"
	youtube "google.golang.org/api/youtube/v3"
)

type Config struct {
	YoutubeClient *youtube.Service
}
type Service struct {
	youtubeClient *youtube.Service
}

func New(config Config) (*Service, error) {
	if config.YoutubeClient == nil {
		return nil, errors.Errorf("%T.APIKey must not be empty", config)
	}

	s := &Service{
		youtubeClient: config.YoutubeClient,
	}

	return s, nil
}

func (s *Service) GetUploadID(youtubeChannel string) (string, error) {
	call := s.youtubeClient.Channels.List("ContentDetails").Id(youtubeChannel)
	response, err := call.Do()
	if err != nil {
		return "", errors.Trace(err)
	}
	if len(response.Items) < 1 {
		return "", errors.New("Not enough results for upload playlist")
	}
	playlistId := response.Items[0].ContentDetails.RelatedPlaylists.Uploads

	return playlistId, nil
}
