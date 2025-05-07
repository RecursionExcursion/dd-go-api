package core

// import (
// 	"errors"
// 	"log"
// 	"strconv"
// )

// type Stat struct {
// 	Value  int
// 	Weight int
// 	Rank   int
// }

// type TrackedStats struct {
// 	Total struct {
// 		Wins         Stat
// 		Losses       Stat
// 		TotalOffense Stat
// 		TotalDefense Stat
// 		PF           Stat
// 		PA           Stat
// 	}
// 	PG struct {
// 		WinsPG   Stat
// 		LossesPG Stat
// 		OffPG    Stat
// 		DefPG    Stat
// 		PFPG     Stat
// 		PAPG     Stat
// 	}
// }

// func (ts *TrackedStats) append(next TrackedStats) {
// 	ts.Total.Wins.Value += next.Total.Wins.Value
// 	ts.Total.Losses.Value += next.Total.Losses.Value

// 	ts.Total.TotalOffense.Value += next.Total.TotalOffense.Value
// 	ts.Total.TotalDefense.Value += next.Total.TotalDefense.Value

// 	ts.Total.PF.Value += next.Total.PF.Value
// 	ts.Total.PA.Value += next.Total.PA.Value
// }

// type WeeklyTeam struct {
// 	Id          int
// 	Week        int
// 	GamesPlayed []string
// 	Stats       TrackedStats
// }

// type ComputedSeason struct {
// 	SeasonInfo         CFBRSeason
// 	RegularSeasonWeeks [][]WeeklyTeam
// 	PostSeasonWeeks    [][]WeeklyTeam
// }

// func ComputeSeason(s CFBRSeason) (ComputedSeason, error) {
// 	cs, err := createWeeks(s)
// 	if err != nil {
// 		return ComputedSeason{}, err
// 	}

// 	//here the cs should be full of each weeks stats

// 	return cs, nil
// }

// func createWeeks(s CFBRSeason) (ComputedSeason, error) {
// 	lw := FindLastWeek(s)

// 	cs := ComputedSeason{
// 		SeasonInfo:         s,
// 		RegularSeasonWeeks: make([][]WeeklyTeam, lw),
// 		PostSeasonWeeks:    make([][]WeeklyTeam, 1),
// 	}

// 	for i := range lw {
// 		currWeek := i + 1
// 		weekTeams := []WeeklyTeam{}

// 		for _, t := range s.Schools {
// 			wt := WeeklyTeam{
// 				Id:          t.Team.Id,
// 				Week:        currWeek,
// 				GamesPlayed: []string{},
// 				Stats:       TrackedStats{},
// 			}

// 			//TODO Get games played
// 			for _, gId := range t.Games {
// 				gm, err := s.FindGameById(gId)
// 				if err != nil {
// 					log.Panicf("Game %v not found", gId)
// 				}

// 				if gm.Game.SeasonType == postseason || gm.Game.Week > currWeek {
// 					continue
// 				}

// 				//Compile stats from game (gm)
// 				ts, err := compileGameStats(t.Team.Id, gm)
// 				if err != nil {
// 					panic(err)
// 				}

// 				wt.Stats.append(ts)
// 			}

// 			//TODO get stats for games played

// 			weekTeams = append(weekTeams, wt)
// 		}
// 		cs.RegularSeasonWeeks[i] = weekTeams
// 	}

// 	return cs, nil
// }

// func FindLastWeek(s CFBRSeason) int {

// 	lastRegSeasonWeek := 0

// 	for _, g := range s.Games {
// 		if g.Game.SeasonType == regularSeason &&
// 			g.Game.Completed &&
// 			g.Game.Week > lastRegSeasonWeek {
// 			lastRegSeasonWeek = g.Game.Week
// 		}
// 	}

// 	return lastRegSeasonWeek
// }

// func compileGameStats(tId int, cg CompleteGame) (TrackedStats, error) {

// 	tm, opp, err := cg.getTeam(tId)
// 	if err != nil {
// 		return TrackedStats{}, err
// 	}

// 	off, err := getTotalYards(tm)
// 	if err != nil {
// 		log.Printf("Off could not be found for team %v in game %v", tId, cg.Id)
// 	}

// 	def, err := getTotalYards(opp)
// 	if err != nil {
// 		log.Printf("Def could not be found for team %v in game %v", tId, cg.Id)
// 	}

// 	win, loss := getWinLoss(tm, opp)
// 	ts := TrackedStats{
// 		Total: struct {
// 			Wins         Stat
// 			Losses       Stat
// 			TotalOffense Stat
// 			TotalDefense Stat
// 			PF           Stat
// 			PA           Stat
// 		}{
// 			Wins:   win,
// 			Losses: loss,

// 			TotalOffense: off,
// 			TotalDefense: def,

// 			PA: getPointsScored(opp),
// 			PF: getPointsScored(tm),
// 		},
// 	}
// 	return ts, nil
// }

// func getTotalYards(tm GameTeam) (Stat, error) {
// 	for _, s := range tm.Stats {
// 		if s.Category == totalYardsStatKey {
// 			stat, err := strconv.Atoi(s.Stat)
// 			if err != nil {

// 			}
// 			return Stat{
// 				Value: stat,
// 			}, nil
// 		}
// 	}
// 	return Stat{}, errors.New("total offense not found")
// }

// func getPointsScored(tm GameTeam) Stat {
// 	return Stat{
// 		Value: int(tm.Points),
// 	}
// }

// func getWinLoss(tm GameTeam, opp GameTeam) (winStat Stat, lossStat Stat) {
// 	if tm.Points > opp.Points {
// 		winStat = Stat{
// 			Value: 1,
// 		}
// 	} else {
// 		lossStat = Stat{
// 			Value: 1,
// 		}
// 	}
// 	return winStat, lossStat
// }
