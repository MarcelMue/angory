package game

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/juju/errors"
)

type Game struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Steam string `json:"steam,omitempty"`
}

func Contains(games []Game, game Game) bool {
	for _, g := range games {
		if g.ID == game.ID {
			return true
		}
	}
	return false
}

func FromPath(path string) ([]Game, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return []Game{}, errors.Trace(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []Game{}, errors.Trace(err)
	}

	var result []Game
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

func ToPath(path string, game Game) error {
	games, err := FromPath(path)
	if err != nil {
		return errors.Trace(err)
	}
	if Contains(games, game) {
		return errors.Errorf("Game with ID %s already exists", game.ID)
	}
	games = append(games, game)

	gamesJSON, err := json.MarshalIndent(games, "", "  ")
	if err != nil {
		return errors.Trace(err)
	}
	err = ioutil.WriteFile(path, gamesJSON, 0644)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
