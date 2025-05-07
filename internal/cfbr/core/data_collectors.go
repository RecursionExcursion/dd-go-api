package core

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

const espnSeasonDateFormat = "2006-01-02T15:04Z"
const espnQueryDateFormat = "20060102"

type SeasonOccurences struct {
	GamesPlayed int
	Schedule    []CollectedGame
}

type CollectedGame struct {
	GameId string
	OppId  string
}
type TeamCollector map[string]*SeasonOccurences

func (tc TeamCollector) Add(c Competitor, opp Competitor, match SeasonCompetition) {
	so, exists := tc[c.Id]

	cg := CollectedGame{
		GameId: match.Id,
		OppId:  opp.Id,
	}

	if exists {
		so.GamesPlayed++
		so.Schedule = append(so.Schedule, cg)
	} else {
		tc[c.Id] = &SeasonOccurences{
			GamesPlayed: 1,
			Schedule:    []CollectedGame{cg},
		}
	}
}

func (tc TeamCollector) FilterFbsTeams() {
	toDelete := []string{}

	for k, v := range tc {
		// most fbs teams play 12+ games, 10 gives it a nice buffer (134 teams in 2024)
		if v.GamesPlayed < 10 {
			toDelete = append(toDelete, k)
		}
	}

	for _, id := range toDelete {
		delete(tc, id)
	}

	/* At this point *most teams will be filtered but.....
	* the geniuses over at ESPN include future fbs addtions
	* so we need to cross ref the scheduled and ensure the majority of games
	* are not paycheck games (fbs vs fcs)
	 */

	toDelete = []string{}
	for k, v := range tc {
		fbsGames := 0
		for _, g := range v.Schedule {
			_, exists := tc[g.OppId]
			if exists {
				fbsGames++
			}
		}
		fbsRatio := float32(fbsGames) / float32(v.GamesPlayed)

		// 50% games are played against fbs teams, this number is negotiable
		if fbsRatio < .5 {
			toDelete = append(toDelete, k)
		}

	}

	for _, id := range toDelete {
		delete(tc, id)
	}
}

func CompileSeason(year int) {
	s, err := getZeroDay(year)
	if err != nil {
		//TODO
		panic(err)
	}

	startDate, endDate, err := getSeasonDateRanges(s)
	if err != nil {
		//TODO
		panic(err)
	}

	collector, err := collectSeason(startDate, endDate)
	if err != nil {
		//TODO
		panic(err)
	}

	// lib.PrettyLog(collector)
	collector.FilterFbsTeams()
	lib.PrettyLog(len(collector))
	lib.PrettyLog(collector["130"])

}

func getZeroDay(year int) (ESPNSeason, error) {
	//0 day query 08/01
	query := fmt.Sprintf("%v0801", year)
	s, err := fetchEspnSeason(query)
	if err != nil {
		return ESPNSeason{}, err
	}
	return s, nil
}

func getSeasonDateRanges(s ESPNSeason) (start time.Time, end time.Time, err error) {
	//get regualr season dates

	if len(s.Leagues) == 0 {
		return time.Time{}, time.Time{}, errors.New("no leagues found")
	}
	for _, c := range s.Leagues[0].Calender {
		if c.Label == "Regular Season" {
			start, err = time.Parse(espnSeasonDateFormat, c.StartDate)
			if err != nil {
				return start, end, err
			}
			end, err = time.Parse(espnSeasonDateFormat, c.EndDate)
			if err != nil {
				return start, end, err
			}
		}
		if c.Label == "Postseason" {
			end, err = time.Parse(espnSeasonDateFormat, c.EndDate)
			if err != nil {
				return start, end, err
			}
		}
	}
	return start, end, nil
}

func collectSeason(startDate time.Time, endDate time.Time) (tc TeamCollector, err error) {
	currDate := startDate

	for {
		//call api
		res, err := fetchEspnSeason(currDate.Format(espnQueryDateFormat))
		if err != nil {
			return tc, err
		}
		//proccess req into map
		for _, e := range res.Events {
			match := e.Competitions[0]
			t1 := match.Competitors[0]
			t2 := match.Competitors[1]

			tc.Add(t1, t2, match)
			tc.Add(t2, t1, match)
		}
		log.Printf("Query for %v complete", currDate)

		//inc date
		currDate = currDate.Add(time.Hour * 24)
		//exit
		if currDate.After(endDate) {
			break
		}
	}

	return tc, nil
}
