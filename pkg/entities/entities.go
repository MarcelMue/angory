package entities

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
	youtube "google.golang.org/api/youtube/v3"
)

type Annotation struct {
	Game    *Game     `json:"game,omitempty"`
	Talents []*Talent `json:"talent,omitempty"`
}

type Game struct {
	Itch  string `json:"itch,omitempty"`
	Name  string `json:"name,omitempty"`
	Steam string `json:"steam,omitempty"`
}

type Statistics struct {
	CommentCount int `json:"commentCount,omitempty"`
	DislikeCount int `json:"dislikeCount,omitempty"`
	LikeCount    int `json:"likeCount,omitempty"`
	ViewCount    int `json:"viewCount,omitempty"`
}

type Talent struct {
	Channel  string `json:"channel,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Short    string `json:"short,omitempty"`
	Twitter  string `json:"twitter,omitempty"`
}

type Video struct {
	Annotation  *Annotation `json:"annotation,omitempty"`
	ChannelID   string      `json:"channelID,omitempty"`
	Description string      `json:"description,omitempty"`
	ID          string      `json:"id"`
	Kind        string      `json:"kind,omitempty"`
	PublishedAt string      `json:"publishedAt,omitempty"`
	Title       string      `json:"titel,omitempty"`
	Statistics  Statistics  `json:"statistics,omitempty"`
	YoutubeTags []string    `json:"youtubeTags,omitempty"`
}

func FromYoutubeVideosPath(path string) ([]*Video, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return []*Video{}, errors.Trace(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []*Video{}, errors.Trace(err)
	}

	var youtubeVideos []*youtube.Video
	json.Unmarshal([]byte(byteValue), &youtubeVideos)

	var result []*Video
	for _, ytVideo := range youtubeVideos {
		vid := &Video{
			ChannelID:   ytVideo.Snippet.ChannelId,
			Description: ytVideo.Snippet.Description,
			ID:          ytVideo.Id,
			Kind:        ytVideo.Kind,
			PublishedAt: ytVideo.Snippet.PublishedAt,
			Title:       ytVideo.Snippet.Title,
			Statistics: Statistics{
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
