package video

type Annotation struct {
	Game   string   `json:"game,omitempty"`
	Talent []string `json:"talent,omitempty"`
}

type Video struct {
	Annotation  *Annotation `json:"annotation,omitempty"`
	ChannelID   string      `json:"channelID,omitempty"`
	Description string      `json:"description,omitempty"`
	ID          string      `json:"id,omitempty"`
	Kind        string      `json:"kind,omitempty"`
	PublishedAt string      `json:"publishedAt,omitempty"`
	Title       string      `json:"titel,omitempty"`
	Statistics  Statistics  `json:"statistics,omitempty"`
	YoutubeTags []string    `json:"youtubeTags,omitempty"`
}

type Statistics struct {
	CommentCount int `json:"commentCount,omitempty"`
	DislikeCount int `json:"dislikeCount,omitempty"`
	LikeCount    int `json:"likeCount,omitempty"`
	ViewCount    int `json:"viewCount,omitempty"`
}
