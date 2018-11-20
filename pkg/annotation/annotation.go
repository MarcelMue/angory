package annotation

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
)

type Annotation struct {
	VideoID   string   `json:"id"`
	GameID    string   `json:"gameID,omitempty"`
	TalentIDs []string `json:"talentIDs,omitempty"`
}

func Contains(annotations []Annotation, annotation Annotation) bool {
	for _, a := range annotations {
		if a.VideoID == annotation.VideoID {
			return true
		}
	}
	return false
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

func ToPath(path string, annotation Annotation) error {
	annotations, err := FromPath(path)
	if err != nil {
		return errors.Trace(err)
	}
	if Contains(annotations, annotation) {
		return errors.Errorf("Annotation with VideoID %s already exists", annotation.VideoID)
	}
	annotations = append(annotations, annotation)

	annotationsJSON, err := json.MarshalIndent(annotations, "", "  ")
	if err != nil {
		return errors.Trace(err)
	}
	err = ioutil.WriteFile(path, annotationsJSON, 0644)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
