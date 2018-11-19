package annotation

type Annotation struct {
	VideoID string   `json:"id,omitempty"`
	Game    string   `json:"game,omitempty"`
	Talent  []string `json:"talent,omitempty"`
}
