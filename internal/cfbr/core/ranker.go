package core

import "log"

type RankerTeam struct {
	Id int
}

type Stat struct {
	Id         int
	TotalYards int
	Points     int
}

type RankerGameStats struct {
	Home Stat
	Away Stat
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

//TODO dont forget PI and SS
type team struct {
	id     int
	week   int
	rank   int
	weight int
	games  []int
	stats  struct {
		Wins         int
		Losses       int
		TotalOffense int
		TotalDefense int
		PF           int
		PA           int
	}
}

type week struct {
	week  int
	games []int
}

type teamMap map[int]RankerTeam
type gameMap map[int]RankerGame
type weekList []week
type weightedWeekMap []map[int]team

func (wkMap *weightedWeekMap) copyWeek(week int, newWeek int) map[int]team {
	mapCopy := map[int]team{}

	for k, v := range (*wkMap)[week] {
		gamesCopy := make([]int, len(v.games))
		copy(gamesCopy, v.games)

		mapCopy[k] = team{
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
func CompileSeasonStats(rs *RankedSeason) *RankedSeason {

	for _, wk := range rs.weeks {

		if len(rs.weightedWeeks) == 0 {
			//Create first week if none exist
			wkMap := make(map[int]team)
			for _, tm := range rs.teams {
				wkMap[tm.Id] = team{
					id:   tm.Id,
					week: wk.week,
				}
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
				UpdateWeightedTeam(homeTeam, awayTeam, &wtHome)
				rs.weightedWeeks[wk.week][homeTeam.Id] = wtHome
			}

			//away team stats
			if wtAway, ok := rs.weightedWeeks[wk.week][awayTeam.Id]; ok {
				UpdateWeightedTeam(awayTeam, homeTeam, &wtAway)
				rs.weightedWeeks[wk.week][awayTeam.Id] = wtAway
			}
		}
	}

	return rs
}

func UpdateWeightedTeam(currTeam Stat, oppTeam Stat, tm *team) {
	if currTeam.Id != tm.id {
		log.Panicf("Invalid team ids (%v-%v)", currTeam.Id, tm.id)
	}

	//Points
	tm.stats.PF += currTeam.Points
	tm.stats.PA += oppTeam.Points

	//Yards
	tm.stats.TotalOffense += currTeam.TotalYards
	tm.stats.TotalDefense += oppTeam.TotalYards

	//W/L
	if currTeam.Points > oppTeam.Points {
		tm.stats.Wins++
	} else {
		tm.stats.Losses++
	}
}
