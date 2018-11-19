package talent

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
)

type Talent struct {
	ID       string `json:"id"`
	Channel  string `json:"channel,omitempty"`
	Name     string `json:"name,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Short    string `json:"short,omitempty"`
	Twitter  string `json:"twitter,omitempty"`
}

func Contains(talents []Talent, talent Talent) bool {
	for _, t := range talents {
		if t.ID == talent.ID {
			return true
		}
	}
	return false
}

func FromPath(path string) ([]Talent, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return []Talent{}, errors.Trace(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []Talent{}, errors.Trace(err)
	}

	var result []Talent
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

func ToPath(path string, talent Talent) error {
	talents, err := FromPath(path)
	if err != nil {
		return errors.Trace(err)
	}
	if Contains(talents, talent) {
		return errors.Errorf("Talent with ID %s already exists", talent.ID)
	}
	talents = append(talents, talent)

	talentsJSON, err := json.MarshalIndent(talents, "", "  ")
	if err != nil {
		return errors.Trace(err)
	}
	err = ioutil.WriteFile(path, talentsJSON, 0644)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
