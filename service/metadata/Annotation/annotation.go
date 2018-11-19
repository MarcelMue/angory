package annotation

type Annotation struct {
	VideoID   string   `json:"id"`
	Game      string   `json:"game,omitempty"`
	TalentIDs []string `json:"talentIDs,omitempty"`
}
