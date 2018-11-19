package talent

type Talent struct {
	ID       string `json:"id"`
	Channel  string `json:"channel,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Short    string `json:"short,omitempty"`
	Twitter  string `json:"twitter,omitempty"`
}
