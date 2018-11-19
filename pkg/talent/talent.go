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
