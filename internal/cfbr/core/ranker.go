package core

import (
	"log"
	"sort"
)

/* =====Types==== */

/* External param types */

/* Top level Ranking DS */
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

type RankedSeason struct {
	teams         teamMap
	games         gameMap
	weeks         weekList
	weightedWeeks weightedWeekMap
}

/* Types for internal calcs */
type WeightedStat struct {
	Rank int
	Val  int
}

type TrackedStats struct {
	Wins         WeightedStat
	Losses       WeightedStat
	TotalOffense WeightedStat
	TotalDefense WeightedStat
	PF           WeightedStat
	PA           WeightedStat
}

// TODO dont forget PI and SS
type team struct {
	id     int
	week   int
	rank   int
	weight int
	games  []int
	stats  TrackedStats
}

type teamMap map[int]RankerTeam
type gameMap map[int]RankerGame
type week struct {
	week  int
	games []int
}
type weekList []week

type wkTeamMap map[int]*team

type TeamAndRank struct {
	Id     int
	Rank   int
	Weight int
}
type weightedWeekMap []wkTeamMap

type rankerParams = struct {
	tms      []*team
	sortFn   (func(i, j int) bool)
	accessor func(*team) int
	assigner func(*team, int)
}

/* Fns */
/* Top Level Ranking Fn */
func Rank(
	teams []RankerTeam,
	games []RankerGame,
) RankedSeason {
	szn := BuildSeason(teams, games)
	szn.CompileSeasonStats()
	szn.CalculateStatRankings()
	return szn
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

/* Ranking Methods */
func (rs *RankedSeason) CompileSeasonStats() {

	var UpdateWeightedTeam = func(currTeam RankerStat, oppTeam RankerStat, tm *team) {
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
				wtHome.games = append(wtHome.games, gmId)
				rs.weightedWeeks[wk.week][homeTeam.Id] = wtHome

			}

			//away team stats
			if wtAway, ok := rs.weightedWeeks[wk.week][awayTeam.Id]; ok {
				UpdateWeightedTeam(awayTeam, homeTeam, wtAway)
				wtAway.games = append(wtAway.games, gmId)
				rs.weightedWeeks[wk.week][awayTeam.Id] = wtAway
			}
		}
	}
}

func (rs *RankedSeason) CalculateStatRankings() {

	//iterate through weightedWeeks and sort stats assign rank
	for _, wk := range rs.weightedWeeks {

		tms := wk.toSlice()

		rankConfigs := []rankerParams{
			//wins
			makeRanker(
				tms,
				func(t *team) int { return t.stats.Wins.Val },
				func(t *team, i int) { t.stats.Wins.Rank = i },
				true,
			),
			// off
			makeRanker(
				tms,
				func(t *team) int { return t.stats.TotalOffense.Val },
				func(t *team, i int) { t.stats.TotalOffense.Rank = i },
				true,
			),
			//pf
			makeRanker(tms,
				func(t *team) int { return t.stats.PF.Val },
				func(t *team, i int) { t.stats.PF.Rank = i },
				true,
			),
			/* Descending stats */

			//losses
			makeRanker(
				tms,
				func(t *team) int { return t.stats.Losses.Val },
				func(t *team, i int) { t.stats.Losses.Rank = i },
				false,
			),
			//def
			makeRanker(
				tms,
				func(t *team) int { return t.stats.TotalDefense.Val },
				func(t *team, i int) { t.stats.TotalDefense.Rank = i },
				false,
			),
			//PA
			makeRanker(
				tms,
				func(t *team) int { return t.stats.PA.Val },
				func(t *team, i int) { t.stats.PA.Rank = i },
				false,
			),
		}

		for _, cfg := range rankConfigs {
			rankStat(cfg)
		}

		//sum weights
		for _, tms := range rs.weightedWeeks {
			for _, tm := range tms {
				tm.sumWeights()
			}
		}

		// Assign overall ranking
		for _, tms := range rs.weightedWeeks {
			tmSlice := tms.toSlice()

			sort.Slice(tmSlice, func(i, j int) bool {
				return tmSlice[i].weight < tmSlice[j].weight
			})

			rankStat(
				makeRanker(
					tmSlice,
					func(t *team) int { return t.weight },
					func(t *team, i int) { t.rank = i },
					false,
				))
		}
	}
}

/* Helper Methods */
func (wk *wkTeamMap) toSlice() []*team {
	tms := make([]*team, len(*wk))
	i := 0
	for _, v := range *wk {
		tms[i] = v
		i++
	}
	return tms
}

func (wk *wkTeamMap) GetRankings() []TeamAndRank {
	tmPtrs := wk.toSlice()
	tms := make([]TeamAndRank, len(tmPtrs))
	for i, t := range tmPtrs {
		tms[i] = TeamAndRank{
			Id:     t.id,
			Rank:   t.rank,
			Weight: t.weight,
		}
	}

	sort.Slice(tms, func(i, j int) bool {
		t1 := tms[i].Rank
		t2 := tms[j].Rank

		if t1 == 0 || t2 == 0 {
			//TODO abstract error to avoid panic
			log.Panicln("Team rank is 0, indicating teams have yet to be ranked")
		}

		return t1 < t2
	})

	return tms
}

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

/* TODO will need to add weight param */
func (tm *team) sumWeights() {
	wt := 0
	wt += tm.stats.Wins.Rank
	wt += tm.stats.Losses.Rank
	wt += tm.stats.TotalOffense.Rank
	wt += tm.stats.TotalDefense.Rank
	wt += tm.stats.PF.Rank
	wt += tm.stats.PA.Rank

	tm.weight = wt
}

/* Fns */

func rankStat(params rankerParams) {
	tms := params.tms
	sort.Slice(tms, params.sortFn)

	rankIndex := 1
	var currVal int

	firstTmFlag := true

	for i, tm := range tms {

		val := params.accessor(tm)
		if firstTmFlag {
			currVal = val
			firstTmFlag = false
		}

		if val != currVal {
			rankIndex = i + 1
			currVal = val
		}
		params.assigner(tm, rankIndex)
	}
}

func makeRanker(
	tms []*team,
	accessor func(*team) int,
	assigner func(*team, int),
	desc bool,
) rankerParams {
	return rankerParams{
		tms:      tms,
		sortFn:   sortTeams(tms, accessor, desc),
		accessor: accessor,
		assigner: assigner,
	}
}

func sortTeams(tms []*team, accessor func(t *team) int, desc bool) func(i, j int) bool {
	return func(i, j int) bool {
		a := accessor(tms[i])
		b := accessor(tms[j])
		if desc {
			return a > b
		}
		return a < b
	}
}
