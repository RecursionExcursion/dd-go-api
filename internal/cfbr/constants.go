package cfbr

const baseRoute = "https://api.collegefootballdata.com"
const teams = "/teams"       //?year=<year>"
const games = "/games"       //?division=<division>&year=<year>&seasonType=<type>" //fbs?
const stats = "/games/teams" //?year=<year>&week=<week>&seasonType=<type>""

// seasonTypes
const regularSeason = "regular"
const postseason = "postseason"

// classifications
const fbs = "fbs"
const fcs = "fcs"
const ii = "ii"
const iii = "iii"

var classes = []string{
	fbs, fcs, ii, iii,
}
