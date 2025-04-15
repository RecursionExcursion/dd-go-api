package cfbr

import "github.com/recursionexcursion/dd-go-api/internal/lib"

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

/* Go batching
 * cfbr only makes 18 req but gets rl pretty quickly, 10 works but is not stable
 * 5 seems good for now
 */
const batchSize = 5

var BatchRunner = lib.RunBatchSizeClosure(batchSize)
