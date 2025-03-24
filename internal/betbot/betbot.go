package betbot

import (
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
		Created: time.Now().Format("01-02-2006T15:04:05"),
		Teams:   teams,
		Games:   games,
	}

	return data, nil
}
