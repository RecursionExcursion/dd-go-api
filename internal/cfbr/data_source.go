package cfbr

import (
	"fmt"
	"net/http"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

type cfbrRoutes = struct {
	teams func(year uint) string
	games func(division string, year uint, seasonType string) string
	stats func(year uint, week uint, seasonType string) string
}

func routeBuilder() cfbrRoutes {
	return cfbrRoutes{
		teams: func(year uint) string {
			return fmt.Sprintf("%v%v?year=%v", baseRoute, teams, year)
		},

		games: func(division string, year uint, seasonType string) string {
			return fmt.Sprintf("%v%v?year=%v&division=%v&seasonType=%v",
				baseRoute,
				games,
				year,
				division,
				seasonType)
		},

		stats: func(year, week uint, seasonType string) string {
			return fmt.Sprintf("%v%v?year=%v&week=%v&seasonType=%v",
				baseRoute,
				stats,
				year,
				week,
				seasonType)
		},
	}
}

func fetchTeams(year uint) ([]Team, error) {
	r := routeBuilder().teams(year)
	return fetchDataToT[[]Team](r)
}

func fetchGames(division string, year uint, seasonType string) ([]Game, error) {
	r := routeBuilder().games(division, year, seasonType)
	return fetchDataToT[[]Game](r)
}

func fetchGameStats(year uint, week uint, seasonType string) ([]GameStats, error) {
	r := routeBuilder().stats(year, week, seasonType)
	return fetchDataToT[[]GameStats](r)
}

func fetchDataToT[T any](route string) (T, error) {
	r, err := reqBuilder(route)
	if err != nil {
		var t T
		return t, err
	}
	t, _, err := lib.FetchAndMap[T](r)
	return t, err
}

func reqBuilder(route string) (func() (*http.Response, error), error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", route, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", lib.EnvGetOrPanic("CFB_API_KEY")))
	return func() (*http.Response, error) {
		return c.Do(req)
	}, nil
}
