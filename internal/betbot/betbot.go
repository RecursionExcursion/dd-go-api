package betbot

import (
	"errors"
	"log"
	"strconv"
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

	FindGame(games, strconv.Itoa(401705613))

	data := FirstShotData{
		Created: time.Now().Format("01-02-2006T15:04:05"),
		Teams:   teams,
		Games:   games,
	}
	FindGameInFsd(data, strconv.Itoa(401705613))

	return data, nil
}

func FindGameInFsd(fsd FirstShotData, id string) (game, error) {
	for _, g := range fsd.Games {
		if g.Id == id {
			log.Printf("%+v", g)
			return g, nil
		}
	}
	var g game
	return g, errors.New("game not found")
}

func FindGame(games []game, id string) (game, error) {
	for _, g := range games {
		if g.Id == id {
			log.Printf("%+v", g)
			return g, nil
		}
	}
	var g game
	return g, errors.New("game not found")
}
