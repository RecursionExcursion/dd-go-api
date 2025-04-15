package cfbr

import (
	"fmt"
	"log"
)

func getNewData(division string, year uint) (CFBRSeason, error) {
	sea := CFBRSeason{}

	//TODO rm logs
	log.Printf("Teams %d", len(teams))
	log.Printf("Games %d", len(games))
	log.Printf("Stats %d", len(stats))
	// log.Println(stats[200])

	teams, games, stats, err := collectDataPoints(year, division)
	if err != nil {
		return sea, err
	}

	return createCfbrTeams(teams, games, stats)
}

func collectDataPoints(year uint, division string) (teams []Team, games []Game, stats []GameStats, err error) {
	teams, err = collectTeams(year, division)
	if err != nil {
		return
	}

	games, err = collectGames(year, division)
	if err != nil {
		return
	}

	/* Team Ids (will be filtered down to div here) */
	tIds := []uint{}
	for _, t := range teams {
		tIds = append(tIds, t.Id)
	}

	stats, err = collectGameStats(year, games, tIds)
	if err != nil {
		return
	}

	return
}

func collectTeams(year uint, division string) ([]Team, error) {
	allTeams, err := fetchTeams(year)
	if err != nil {
		return nil, err
	}

	//TODO move to filter mod???
	/* Filter by division */
	divisionTeams := []Team{}

	for _, t := range allTeams {
		if t.Classification == division {
			divisionTeams = append(divisionTeams, t)
		}
	}

	return divisionTeams, nil
}

func collectGames(year uint, division string) ([]Game, error) {
	games := []Game{}

	gReg, err := fetchGames(division, year, "regular")
	if err != nil {
		return nil, err
	}

	gPost, err := fetchGames(division, year, "postseason")
	if err != nil {
		return nil, err
	}

	addGames := func(gms []Game) {
		for _, g := range gms {
			if g.Completed {
				games = append(games, g)
			}
		}
	}

	addGames(gReg)
	addGames(gPost)

	return games, nil
}

func collectGameStats(year uint, games []Game, teamIds []uint) ([]GameStats, error) {
	regSeason := []Game{}
	postSeason := []Game{}

	for _, g := range games {
		if g.SeasonType == regularSeason {
			regSeason = append(regSeason, g)
		} else if g.SeasonType == postseason {
			postSeason = append(postSeason, g)
		} else {
			log.Printf("No season type found on game %v", g.Id)
		}

	}

	//Calc max week
	maxWeek := 0
	for _, g := range regSeason {
		if g.Completed && g.Week > uint(maxWeek) {
			maxWeek = int(g.Week)
		}
	}

	allGameStats := []GameStats{}

	// regular season
	for i := 0; i <= maxWeek; i++ {
		gs, err := fetchGameStats(year, uint(i), "regular")
		if err != nil {
			return nil, err
		}
		allGameStats = append(allGameStats, gs...)
	}

	for _, g := range postSeason {
		if g.Completed && g.Week > uint(maxWeek) {
			maxWeek = int(g.Week)
		}
	}

	//postseason
	gs, err := fetchGameStats(year, 1, "postseason")
	if err != nil {
		return nil, err
	}
	allGameStats = append(allGameStats, gs...)

	//TODO move to filter mod???
	//Filter gamestats against div

	filteredGs := []GameStats{}

	for _, st := range allGameStats {
		t1 := st.Teams[0]
		t2 := st.Teams[1]

		for _, tId := range teamIds {
			if tId == t1.SchoolId || tId == t2.SchoolId {
				filteredGs = append(filteredGs, st)
				break
			}
		}

	}

	log.Printf("All %v", len(allGameStats))
	log.Printf("Filt %v", len(filteredGs))

	return filteredGs, nil
}

func getGameStatsById(stats []GameStats, id uint) (GameStats, error) {
	for _, gs := range stats {

		if gs.Id == id {
			return gs, nil
		}

	}

	return GameStats{}, fmt.Errorf("could not find game stat %v", id)
}

func createCfbrTeams(teams []Team, games []Game, stats []GameStats) (CFBRSeason, error) {
	sea := EmptySeason()

	for _, t := range teams {
		sea.Schools[t.Id] = CFBRSchool{
			Team:  t,
			Games: []CompleteGame{},
		}
	}

	for _, g := range games {

		gs, err := getGameStatsById(stats, g.Id)
		if err != nil {
			return sea, err
		}

		homeSchool, ok := sea.Schools[g.HomeId]
		if ok {
			homeSchool.Games = append(homeSchool.Games,
				CompleteGame{
					Game:      g,
					GameStats: gs,
				},
			)

			sea.Schools[g.HomeId] = homeSchool
		}

		awaySchool, ok := sea.Schools[g.AwayId]
		if ok {
			awaySchool.Games = append(awaySchool.Games,
				CompleteGame{
					Game:      g,
					GameStats: gs,
				},
			)

			sea.Schools[g.AwayId] = awaySchool
		}

	}

	return sea, nil
}
