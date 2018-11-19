package video

import (
	"sort"
	"strings"

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

func (s *Service) FromPlaylist(playlistID string) ([]*youtube.Video, error) {
	var allVideos []*youtube.Video

	{
		nextPageToken := ""
		for {
			call := s.youtubeClient.PlaylistItems.List("ContentDetails,Snippet").PlaylistId(playlistID).MaxResults(50).PageToken(nextPageToken)
			response, err := call.Do()
			if err != nil {
				return nil, errors.Trace(err)
			}
			var pageVideoIDs string
			{
				allSlice := []string{}
				for _, item := range response.Items {
					allSlice = append(allSlice, item.ContentDetails.VideoId)
				}
				pageVideoIDs = strings.Join(allSlice, ",")
			}
			pageVideos, err := s.FromIDs(pageVideoIDs)
			if err != nil {
				return nil, errors.Trace(err)
			}
			allVideos = append(allVideos, pageVideos...)

			sort.Slice(allVideos, func(i int, j int) bool {
				return allVideos[i].Snippet.PublishedAt > allVideos[j].Snippet.PublishedAt
			})

			nextPageToken = response.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}

	return allVideos, nil
}

func (s *Service) FromIDs(videoIDs string) ([]*youtube.Video, error) {
	var allVideos []*youtube.Video
	{
		nextPageToken := ""
		for {
			call := s.youtubeClient.Videos.List("Snippet,Statistics").Id(videoIDs).MaxResults(50).PageToken(nextPageToken)
			response, err := call.Do()
			if err != nil {
				return nil, errors.Trace(err)
			}
			allVideos = append(allVideos, response.Items...)
			nextPageToken = response.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}

	return allVideos, nil
}
