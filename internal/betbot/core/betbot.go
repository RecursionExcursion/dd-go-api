package core

import (
	"errors"
	"time"
)

func CollectData() (FirstShotData, error) {

	teams, err := collectTeamsAndRosters()
	if err != nil {
		return FirstShotData{}, err
	}
	games, err := collectGames()
	if err != nil {
		return FirstShotData{}, err
	}

	data := FirstShotData{
		Created: time.Now().UnixMilli(),
		Teams:   teams,
		Games:   games,
	}

	return data, nil
}

func FindGameInFsd(fsd FirstShotData, id string) (game, error) {
	for _, g := range fsd.Games {
		if g.Id == id {
			return g, nil
		}
	}
	var g game
	return g, errors.New("game not found")
}

func FindGame(games []game, id string) (game, error) {
	for _, g := range games {
		if g.Id == id {
			return g, nil
		}
	}
	var g game
	return g, errors.New("game not found")
}
