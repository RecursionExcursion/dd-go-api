package core

// import (
// 	"errors"
// 	"fmt"
// )

// type CompleteGame struct {
// 	Id        int
// 	Game      Game
// 	GameStats GameStats
// }

// func (cg *CompleteGame) getTeam(id int) (currTeam GameTeam, oppTeam GameTeam, err error) {

// 	oppAssignedFlag := false

// 	for _, t := range cg.GameStats.Teams {
// 		if t.SchoolId == id {
// 			currTeam = t
// 		} else {
// 			if oppAssignedFlag == true {
// 				err = fmt.Errorf("Team %v not found", id)
// 			}
// 			oppTeam = t
// 			oppAssignedFlag = true
// 		}
// 	}

// 	return currTeam, oppTeam, err
// }

// type CFBRSchool struct {
// 	Team  Team
// 	Games []int
// }

// /* CFBRSeason- This is the main data structure for this module
// *
//  */

type SerializeableCompressedSeason struct {
	Id               string `json:"id"`
	CreatedAt        int    `json:"createdAt"`
	Year             int    `json:"year"`
	CompressedSeason string `json:"season"`
}

// type GameMap = map[string]CompleteGame
// type SchoolMap = map[string]CFBRSchool

// type CFBRSeason struct {
// 	Year     int
// 	Division string
// 	Schools  SchoolMap
// 	Games    GameMap
// }

// func EmptySeason() CFBRSeason {
// 	return CFBRSeason{
// 		Schools: make(SchoolMap),
// 	}
// }

// // First accept args (stat weights)/(year)
// func Create(divsion string, year int) (CFBRSeason, error) {
// 	// season, err := collectCfbSeasonData(divsion, year)
// 	// if err != nil {
// 	// 	return CFBRSeason{}, err
// 	// }
// 	// return season, nil
// 	return CFBRSeason{}, nil
// }

// /* CFBRSeason- Util fns */

// func (c *CFBRSeason) FindSchoolById(id int) (s CFBRSchool, err error) {
// 	s, ok := c.Schools[fmt.Sprint(id)]
// 	if !ok {
// 		err = errors.New("school not found")
// 	}
// 	return s, err
// }

// func (c *CFBRSeason) FindGameById(id int) (cg CompleteGame, err error) {
// 	cg, ok := c.Games[fmt.Sprint(id)]
// 	if !ok {
// 		err = errors.New("game not found")
// 	}
// 	return cg, err
// }
