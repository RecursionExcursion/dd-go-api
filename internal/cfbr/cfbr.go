package cfbr

import (
	"errors"
	"fmt"
	"strconv"
)

/* CFBRSeason- This is the main data structure for this module
 *
 */

type CFBRSeason struct {
	Schools map[uint]CFBRSchool
}

func EmptySeason() CFBRSeason {
	return CFBRSeason{
		Schools: make(map[uint]CFBRSchool),
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

func (c *CFBRSeason) Save() map[string]CFBRSchool {
	outMap := make(map[string]CFBRSchool, len(c.Schools))

	for k, v := range c.Schools {
		outMap[fmt.Sprint(k)] = v
	}
	return outMap
}

func (c *CFBRSeason) Load(inMap map[string]CFBRSchool) (CFBRSeason, error) {
	schoolMap := make(map[uint]CFBRSchool, len(inMap))

	for k, v := range inMap {
		n, err := strconv.ParseUint(k, 10, 0)
		if err != nil {
			return CFBRSeason{}, err
		}

		schoolMap[uint(n)] = v
	}

	return CFBRSeason{
		Schools: schoolMap,
	}, nil

}

/* CFBRSeason- Util fns */

func (c *CFBRSeason) FindSchoolById(id uint) (s CFBRSchool, err error) {
	s, ok := c.Schools[id]
	if !ok {
		err = errors.New("school not found")
	}

	return s, err
}

func (c *CFBRSeason) FindGameById(id uint) (CompleteGame, error) {

	for _, s := range c.Schools {
		g, err := s.FindGameById(id)
		if err != nil {
			continue
		}
		return g, nil
	}

	return CompleteGame{}, errors.New("Team not found")
}

/* Compiled School */

type CompleteGame struct {
	Game      Game
	GameStats GameStats
}

type CFBRSchool struct {
	Team  Team
	Games []CompleteGame
}

func (s *CFBRSchool) FindGameById(id uint) (CompleteGame, error) {
	for _, g := range s.Games {
		if g.Game.Id == id {
			return g, nil
		}
	}

	return CompleteGame{}, errors.New("Game not found")
}
