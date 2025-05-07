package core

import (
	"fmt"
	"log"
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

type cfbrRoutes = struct {
	// cfbApi cfbApi
	espn espnApi
}

// type cfbApi struct {
// 	teams func(year int) string
// 	games func(division string, year int, seasonType string) string
// 	stats func(year int, week int, seasonType string) string
// }

type espnApi struct {
	groups   func() string
	season   func(year string) string
	allTeams func() string
	team     func(teamId int) string
	stats    func(eventId int) string
	day      func(date string) string
}

func routeBuilder() cfbrRoutes {
	return cfbrRoutes{
		// cfbApi: cfbApi{
		// 	teams: func(year int) string {
		// 		return fmt.Sprintf("%v%v?year=%v", baseRoute, teams, year)
		// 	},

		// 	games: func(division string, year int, seasonType string) string {
		// 		return fmt.Sprintf("%v%v?year=%v&division=%v&seasonType=%v",
		// 			baseRoute,
		// 			games,
		// 			year,
		// 			division,
		// 			seasonType)
		// 	},

		// 	stats: func(year, week int, seasonType string) string {
		// 		return fmt.Sprintf("%v%v?year=%v&week=%v&seasonType=%v",
		// 			baseRoute,
		// 			stats,
		// 			year,
		// 			week,
		// 			seasonType)
		// 	},
		// },
		espn: espnApi{
			groups: func() string {
				return fmt.Sprintf("%v%v", espnBase, espnGroups)
			},
			season: func(date string) string {
				return fmt.Sprintf("%v%v?dates=%v", espnBase, espnSeason, date)
			},
			allTeams: func() string {
				return fmt.Sprintf("%v%v", espnBase, espnTeams)
			},
			team: func(teamId int) string {
				return fmt.Sprintf("%v%v/%v", espnBase, espnTeams, teamId)
			},
			stats: func(eventId int) string {
				return fmt.Sprintf("%v%v?event=%v", espnBase, espnGame, eventId)
			},
		},
	}
}

func fetchEspnAllTeams() (ESPNTeams, error) {
	r := routeBuilder().espn.allTeams()
	return fetchDataToT[ESPNTeams](r)
}

func fetchEspnTeam(teamId int) (ESPNCfbTeam, error) {
	r := routeBuilder().espn.team(teamId)
	return fetchDataToT[ESPNCfbTeam](r)
}

// func fetchEspnGroups() (ESPNGroups, error) {
// }

func fetchEspnSeason(date string) (ESPNSeason, error) {
	r := routeBuilder().espn.season(date)
	return fetchDataToT[ESPNSeason](r)
}

func fetchEspnStats(eventId int) (ESPNCfbGame, error) {
	r := routeBuilder().espn.stats(eventId)
	return fetchDataToT[ESPNCfbGame](r)
}

// func fetchTeams(year int) ([]Team, error) {
// 	r := routeBuilder().teams(year)
// 	return fetchDataToT[[]Team](r)
// }

// func fetchGames(division string, year int, seasonType string) ([]Game, error) {
// 	r := routeBuilder().games(division, year, seasonType)
// 	return fetchDataToT[[]Game](r)
// }

// func fetchGameStats(year int, week int, seasonType string) ([]GameStats, error) {
// 	r := routeBuilder().stats(year, week, seasonType)
// 	return fetchDataToT[[]GameStats](r)
// }

func fetchDataToT[T any](route string) (T, error) {
	r, err := reqBuilder(route)
	if err != nil {
		var t T
		return t, err
	}
	t, res, err := lib.FetchAndMap[T](r)
	if err != nil {
		log.Printf("Request to %v failed", res.Request.URL)
	}
	return t, err
}

func reqBuilder(route string) (func() (*http.Response, error), error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", lib.EnvGetOrPanic("CFB_API_KEY")))
	return func() (*http.Response, error) {
		return c.Do(req)
	}, nil
}
