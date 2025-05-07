package core

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

/* CFBR batching
 * cfbr only makes 18 req but gets ratelimited pretty quickly, 10 works but is not stable (yet?)
 * 5 seems safe for now until a more robust ratelmiting logic is impl
 */
const batchSize = 5

var BatchRunner = lib.RunBatchSizeClosure(batchSize)

var trackedStatCategories = []string{
	"totalYards",
}

var totalYardsStatKey = "totalYards"

/* ESPN Routes */
const espnBase = "https://site.api.espn.com/apis/site/v2/sports/football/college-football"
const espnGroups = "/groups"
const espnSeason = "/scoreboard" //dates=2024 or dates=20240921
const espnTeams = "/teams"       //</teamid>
const espnGame = "/summary"      //?event=<eventId>

/* Group Keys */
type GroupName = struct {
	name     string
	children []string
}

var D1 = GroupName{
	name: "NCAA Division I",
	children: []string{
		"FBS (I-A)",
		"FCS (I-AA)",
	},
}

var D2_3 = GroupName{
	name: "Division II/III",
	children: []string{
		"NCAA Division II",
		"NCAA Division III",
	},
}

/* TODO: DELETE place holder for flow
 * Groups (Teams) -> Season (Games) -> Games
 * The actual teams endpoint may be  moot but we will see what we need from it.
 *
 *
 */
