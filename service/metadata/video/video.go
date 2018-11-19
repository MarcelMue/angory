package video

type Annotation struct {
	Game    string   `json:"game,omitempty"`
	Talents []Talent `json:"talent,omitempty"`
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
