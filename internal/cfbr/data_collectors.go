package cfbr

import (
	"fmt"
	"log"
)

func collectCfbSeasonData(division string, year uint) (CFBRSeason, error) {
	sea := CFBRSeason{}

	//TODO add concurretcy
	teams, games, stats, err := collectDataPoints(year, division)
	if err != nil {
		return sea, err
	}

	//TODO rm logs
	log.Printf("Teams %d", len(teams))
	log.Printf("Games %d", len(games))
	log.Printf("Stats %d", len(stats))
	// log.Println(stats[200])

	return createCfbrTeams(teams, games, stats)
}

func collectDataPoints(year uint, division string) (teams []Team, games []Game, stats []GameStats, err error) {

	tChan := make(chan []Team)
	gChan := make(chan []Game)
	tasks := []func(){
		func() {
			teams, err = collectTeams(year, division)
			if err != nil {
				log.Println(err)
				tChan <- []Team{}
				return
			}
			tChan <- teams
		},
		func() {
			games, err = collectGames(year, division)
			if err != nil {
				log.Println(err)
				gChan <- []Game{}
				return
			}
			gChan <- games
		},
	}

	go func() {
		BatchRunner(tasks)
	}()

	teams = <-tChan
	games = <-gChan

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

	gChan := make(chan []Game)
	tasks := []func(){
		func() {
			gReg, err := fetchGames(division, year, "regular")
			if err != nil {
				gChan <- []Game{}
				return
			}
			gChan <- gReg
		},
		func() {
			gPost, err := fetchGames(division, year, "postseason")
			if err != nil {
				gChan <- []Game{}
				return
			}
			gChan <- gPost
		},
	}

	go func() {
		BatchRunner(tasks)
		close(gChan)
	}()

	games := []Game{}
	for chGames := range gChan {

		for _, g := range chGames {
			if g.Completed {
				games = append(games, g)
			}
		}
	}

	return games, nil
}

func collectGameStats(year uint, games []Game, teamIds []uint) ([]GameStats, error) {
	allGameStats := []GameStats{}

	gsChan := make(chan []GameStats)
	tasks := []func(){}

	//Calc max week for reg season
	maxWeek := 0
	for _, g := range games {
		if g.Completed && g.SeasonType == regularSeason && g.Week > uint(maxWeek) {
			maxWeek = int(g.Week)
		}
	}

	// regular season
	for i := 0; i <= maxWeek; i++ {

		tasks = append(tasks, func() {
			gs, err := fetchGameStats(year, uint(i), "regular")
			if err != nil {
				log.Println(err)
				gsChan <- []GameStats{}
				return
			}
			gsChan <- gs
		})
	}

	//postseason
	tasks = append(tasks, func() {

		gs, err := fetchGameStats(year, 1, "postseason")
		if err != nil {
			log.Println(err)
			gsChan <- []GameStats{}
			return
		}
		gsChan <- gs
	})

	go func() {
		BatchRunner(tasks)
		close(gsChan)
	}()

	for gs := range gsChan {
		allGameStats = append(allGameStats, gs...)
	}

	//TODO move to filter mod???
	//Filter gamestats against div

	log.Printf("Allgamestats: %v", len(allGameStats))

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
