package core

import (
	"log"
	"sort"
)

type RankerTeam struct {
	Id int
}

type RankerStat struct {
	Id         int
	TotalYards int
	Points     int
}

type RankerGameStats struct {
	Home RankerStat
	Away RankerStat
}

type RankerGame struct {
	Id    int
	Week  int
	Stats RankerGameStats
}

//

/* Methods */
func Rank(
	teams []RankerTeam,
	games []RankerGame,
) {
	_ = BuildSeason(teams, games)
}

type WeightedStat struct {
	Rank int
	Val  int
}

// TODO dont forget PI and SS
type team struct {
	id     int
	week   int
	rank   int
	weight int
	games  []int
	stats  struct {
		Wins         WeightedStat
		Losses       WeightedStat
		TotalOffense WeightedStat
		TotalDefense WeightedStat
		PF           WeightedStat
		PA           WeightedStat
	}
}

type week struct {
	week  int
	games []int
}

type teamMap map[int]RankerTeam
type gameMap map[int]RankerGame
type weekList []week
type weightedWeekMap []map[int]*team

func (wkMap *weightedWeekMap) copyWeek(week int, newWeek int) map[int]*team {
	mapCopy := map[int]*team{}

	for k, v := range (*wkMap)[week] {
		gamesCopy := make([]int, len(v.games))
		copy(gamesCopy, v.games)

		mapCopy[k] = &team{
			id:     v.id,
			week:   newWeek,
			rank:   v.rank,
			weight: v.weight,
			games:  gamesCopy,
			stats:  v.stats,
		}
	}

	return mapCopy
}

type RankedSeason struct {
	teams         teamMap
	games         gameMap
	weeks         weekList
	weightedWeeks weightedWeekMap
}

func BuildSeason(teams []RankerTeam, games []RankerGame) RankedSeason {

	tm := teamMap{}
	for _, t := range teams {
		tm[t.Id] = t
	}

	gm := gameMap{}
	wl := weekList{}
	for _, g := range games {
		gm[g.Id] = g

		if g.Week+1 > len(wl) {
			tmp := make(weekList, g.Week+1)
			copy(tmp, wl)
			wl = tmp
			wl[g.Week] = week{
				week:  g.Week,
				games: []int{},
			}
		}
		wk := wl[g.Week]
		wk.games = append(wk.games, g.Id)
		wl[g.Week] = wk
	}

	weightMap := weightedWeekMap{}
	// for _, wk := range wl {
	// 	//Copy all teams to a map
	// 	wkMap := make(map[int]team)
	// 	for _, tm := range teams {
	// 		wkMap[tm.Id] = team{
	// 			id:   tm.Id,
	// 			week: wk.week,
	// 		}
	// 	}

	// 	weightMap = append(weightMap, wkMap)
	// }

	return RankedSeason{
		teams:         tm,
		games:         gm,
		weeks:         wl,
		weightedWeeks: weightMap,
	}
}

// Also will take in ranking params
func CompileSeasonStats(rs *RankedSeason) {

	for _, wk := range rs.weeks {

		if len(rs.weightedWeeks) == 0 {
			//Create first week if none exist
			wkMap := make(map[int]*team)
			for _, tm := range rs.teams {
				newTm := team{
					id:   tm.Id,
					week: wk.week,
				}
				wkMap[tm.Id] = &newTm
			}
			rs.weightedWeeks = append(rs.weightedWeeks, wkMap)
		} else {
			//copy prev week into curr week
			cpy := rs.weightedWeeks.copyWeek(wk.week-1, wk.week)
			rs.weightedWeeks = append(rs.weightedWeeks, cpy)
		}

		for _, gmId := range wk.games {

			gm, ok := rs.games[gmId]
			if !ok {
				//TODO
				log.Panicf("Game id (%v) not found", gmId)
			}

			homeTeam := gm.Stats.Home
			awayTeam := gm.Stats.Away

			//Home team stats
			if wtHome, ok := rs.weightedWeeks[wk.week][homeTeam.Id]; ok {
				UpdateWeightedTeam(homeTeam, awayTeam, wtHome)
				rs.weightedWeeks[wk.week][homeTeam.Id] = wtHome
			}

			//away team stats
			if wtAway, ok := rs.weightedWeeks[wk.week][awayTeam.Id]; ok {
				UpdateWeightedTeam(awayTeam, homeTeam, wtAway)
				rs.weightedWeeks[wk.week][awayTeam.Id] = wtAway
			}
		}
	}
}

func UpdateWeightedTeam(currTeam RankerStat, oppTeam RankerStat, tm *team) {
	if currTeam.Id != tm.id {
		log.Panicf("Invalid team ids (%v-%v)", currTeam.Id, tm.id)
	}

	//Points
	tm.stats.PF.Val += currTeam.Points
	tm.stats.PA.Val += oppTeam.Points

	//Yards
	tm.stats.TotalOffense.Val += currTeam.TotalYards
	tm.stats.TotalDefense.Val += oppTeam.TotalYards

	//W/L
	if currTeam.Points > oppTeam.Points {
		tm.stats.Wins.Val++
	} else {
		tm.stats.Losses.Val++
	}
}

/*
Sort stats by rank in sep slices (this needs to be held in its own ds, and returned)
//TODO HANDLE TIED STATS (will req a fn)
*/

func CalculateStatRankings(rs *RankedSeason) {

	//iterate through weightedWeeks and sort stats assign rank
	for _, wk := range rs.weightedWeeks {
		tms := make([]*team, len(wk))
		for _, v := range wk {
			tms = append(tms, v)
		}

		type rankerParams = struct {
			sortFn       (func(i, j int) bool)
			statAccessor func(*team) *WeightedStat
		}

		rankStat := func(params rankerParams) {
			sort.Slice(tms, params.sortFn)

			rankIndex := 1
			var currVal int

			firstTmFlag := true

			for i, tm := range tms {

				stat := params.statAccessor(tm)
				val := stat.Val
				if firstTmFlag {
					currVal = val
					firstTmFlag = false
				}

				if val != currVal {
					rankIndex = i + 1
					currVal = val
				}
				stat.Rank = rankIndex
			}
		}
		//Sort by stats and assign rank

		//wins
		rankStat(rankerParams{
			sortFn:       func(a, b int) bool { return tms[a].stats.Wins.Val > tms[b].stats.Wins.Val },
			statAccessor: func(t *team) *WeightedStat { return &t.stats.Wins },
		})

		// //off
		rankStat(rankerParams{
			sortFn:       func(a, b int) bool { return tms[a].stats.TotalOffense.Val > tms[b].stats.TotalOffense.Val },
			statAccessor: func(t *team) *WeightedStat { return &t.stats.TotalOffense },
		})

		//pf
		rankStat(rankerParams{
			sortFn:       func(a, b int) bool { return tms[a].stats.PF.Val > tms[b].stats.PF.Val },
			statAccessor: func(t *team) *WeightedStat { return &t.stats.PF },
		})

		/* Descending stats */

		//losess
		rankStat(rankerParams{
			sortFn:       func(a, b int) bool { return tms[a].stats.Losses.Val < tms[b].stats.Losses.Val },
			statAccessor: func(t *team) *WeightedStat { return &t.stats.Losses },
		})

		//def
		rankStat(rankerParams{
			sortFn:       func(a, b int) bool { return tms[a].stats.TotalDefense.Val < tms[b].stats.TotalDefense.Val },
			statAccessor: func(t *team) *WeightedStat { return &t.stats.TotalDefense },
		})

		//PA
		rankStat(rankerParams{
			sortFn:       func(a, b int) bool { return tms[a].stats.PA.Val < tms[b].stats.PA.Val },
			statAccessor: func(t *team) *WeightedStat { return &t.stats.PA },
		})
	}
}
