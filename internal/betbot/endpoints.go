package betbot

import "strings"

const base = "https://site.api.espn.com/apis/site/v2/sports/basketball/nba"

const season = base + "/scoreboard?year={year}"
const scoreboard = base + "/scoreboard?dates={date}"
const gameData = base + "/summary?event={gameID}"
const teams = base + "/teams"
const roster = base + "/teams/{teamID}/roster"

/* ESPN NBA API EndPoint accessors */
type endPointsShape struct {
	Season     func(string) string
	Scoreboard func(string) string
	GameData   func(string) string
	Teams      func() string
	Roster     func(string) string
}

func endpoints() endPointsShape {

	endpoints := endPointsShape{
		Season: func(s string) string {
			return delimit(season, "{year}", s)
		},
		Scoreboard: func(s string) string {
			return delimit(scoreboard, "{date}", s)
		},
		GameData: func(s string) string {
			return delimit(gameData, "{gameID}", s)
		},
		Teams: func() string {
			return delimit(teams, "", "")
		},
		Roster: func(s string) string {
			return delimit(roster, "{teamID}", s)
		},
	}

	return endpoints

}

func delimit(base string, delimitter string, arg string) string {
	return strings.Replace(base, delimitter, arg, 1)
}
