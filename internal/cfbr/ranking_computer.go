package cfbr

import (
	"fmt"
	"log"
	"strconv"
)

type Stat struct {
	Value  int
	Weight int
	Rank   int
}

type TrackedStats struct {
	Total struct {
		Wins         Stat
		Losses       Stat
		TotalOffense Stat
		TotalDefense Stat
		PF           Stat
		PA           Stat
	}
	PG struct {
		WinsPG   Stat
		LossesPG Stat
		OffPG    Stat
		DefPG    Stat
		PFPG     Stat
		PAPG     Stat
	}
}

type WeeklyTeam struct {
	Id          uint
	Week        int
	GamesPlayed []string
	Stats       TrackedStats
}

type ComputedSeason struct {
	SeasonInfo         CFBRSeason
	RegularSeasonWeeks [][]WeeklyTeam
	PostSeasonWeeks    [][]WeeklyTeam
}

func ComputeSeason(s CFBRSeason) (ComputedSeason, error) {
	lw := FindLastWeek(s)

	cs := ComputedSeason{
		SeasonInfo:         s,
		RegularSeasonWeeks: make([][]WeeklyTeam, lw),
		PostSeasonWeeks:    make([][]WeeklyTeam, 1),
	}

	//DO ALOT OF WORK
	for i := range lw {
		currWeek := i + 1
		weekTeams := make([]WeeklyTeam, len(s.Schools))

		for _, t := range s.Schools {
			wt := WeeklyTeam{
				Id:          t.Team.Id,
				Week:        currWeek,
				GamesPlayed: []string{},
				Stats:       TrackedStats{},
			}

			//TODO Get games played
			for _, gId := range t.Games {
				gm, err := s.FindGameById(gId)
				if err != nil {
					log.Panicf("Game %v not found", gId)
				}

				//Compile stats from game (gm)
				ts, err := compileGameStats(t.Team.Id, gm)
			}

			//TODO get stats for games played

			weekTeams = append(weekTeams, wt)
		}
		cs.RegularSeasonWeeks = append(cs.RegularSeasonWeeks, weekTeams)
	}

	return cs, nil
}

func FindLastWeek(s CFBRSeason) int {

	lastRegSeasonWeek := 0

	for _, g := range s.Games {
		if g.Game.SeasonType == regularSeason &&
			g.Game.Completed &&
			g.Game.Week > uint(lastRegSeasonWeek) {
			lastRegSeasonWeek = int(g.Game.Week)
		}
	}

	return lastRegSeasonWeek
}

func compileGameStats(teamId uint, game CompleteGame) (TrackedStats, error) {

	tm, opp, err := func() (GameTeam, GameTeam, error) {
		currTeam := GameTeam{}
		oppTeam := GameTeam{}

		for _, t := range game.GameStats.Teams {
			if t.SchoolId == teamId {
				currTeam = t
			}
		}

		return currTeam, oppTeam, fmt.Errorf("Team %v not found", teamId)
	}()
	if err != nil {
		return TrackedStats{}, err
	}

	ts := TrackedStats{}

	//OFF
	for _, s := range tm.Stats {
		if s.Category == totalYardsStatKey {
			stat, err := strconv.Atoi(s.Stat)
			if err != nil {

			}
			ts.Total.TotalOffense = Stat{
				Value: stat,
			}
		}
	}

	//DEF
	for _, s := range opp.Stats {
		if s.Category == totalYardsStatKey {
			stat, err := strconv.Atoi(s.Stat)
			if err != nil {

			}
			ts.Total.TotalDefense = Stat{
				Value: stat,
			}
		}
	}

	//PA
	ts.Total.PA = Stat{
		Value: int(opp.Points),
	}

	//PF
	ts.Total.PF = Stat{
		Value: int(tm.Points),
	}

	//W/L
	if tm.Points > opp.Points {
		ts.Total.Wins = Stat{
			Value: 1,
		}
	} else {
		ts.Total.Losses = Stat{
			Value: 1,
		}
	}

	return ts, nil
}
