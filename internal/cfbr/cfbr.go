package cfbr

import (
	"errors"
	"fmt"
)

/* Compiled School */

type CompleteGame struct {
	Id        uint
	Game      Game
	GameStats GameStats
}

type CFBRSchool struct {
	Team  Team
	Games []uint
}

/* CFBRSeason- This is the main data structure for this module
*
 */

type SerializeableCompressedSeason struct {
	Id               string `json:"id"`
	CreatedAt        int    `json:"createdAt"`
	Year             int    `json:"year"`
	CompressedSeason string `json:"season"`
}

type GameMap = map[string]CompleteGame
type SchoolMap = map[string]CFBRSchool

type CFBRSeason struct {
	Year     int
	Division string
	Schools  SchoolMap
	Games    GameMap
}

func EmptySeason() CFBRSeason {
	return CFBRSeason{
		Schools: make(SchoolMap),
	}
}

// First accept args (stat weights)/(year)
func Create(divsion string, year uint) (CFBRSeason, error) {
	season, err := collectCfbSeasonData(divsion, year)
	if err != nil {
		return CFBRSeason{}, err
	}
	return season, nil
}

/* CFBRSeason- Util fns */

func (c *CFBRSeason) FindSchoolById(id uint) (s CFBRSchool, err error) {
	s, ok := c.Schools[fmt.Sprint(id)]
	if !ok {
		err = errors.New("school not found")
	}
	return s, err
}

func (c *CFBRSeason) FindGameById(id uint) (cg CompleteGame, err error) {
	cg, ok := c.Games[fmt.Sprint(id)]
	if !ok {
		err = errors.New("game not found")
	}
	return cg, err
}
