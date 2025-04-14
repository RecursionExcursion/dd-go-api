package cfbr

import "errors"

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

/* Season */

type CFBRSeason struct {
	Schools map[uint]CFBRSchool
}

func newCFBRSeason() CFBRSeason {
	return CFBRSeason{
		Schools: make(map[uint]CFBRSchool),
	}
}

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
