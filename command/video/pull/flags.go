package pull

import (
	"github.com/juju/errors"
)

type Flags struct {
	APIKey         string
	YoutubeChannel string
	Path           string
}

func (f *Flags) Validate() error {
	if f.APIKey == "" {
		return errors.Errorf("Flag for API key must not be empty")
	}
	return nil
}
