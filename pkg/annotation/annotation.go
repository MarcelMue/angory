package annotation

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
)

type Annotation struct {
	VideoID   string   `json:"id"`
	Game      string   `json:"game,omitempty"`
	TalentIDs []string `json:"talentIDs,omitempty"`
}

func FromPath(path string) ([]Annotation, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return []Annotation{}, errors.Trace(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []Annotation{}, errors.Trace(err)
	}

	var result []Annotation
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}
