package core

import (
	"errors"
	"fmt"
	"log"
	"sort"
)

//TODO Write test for postseason, and mapping

/* =====Types==== */

/* External param types */

/* Top level Ranking DS */
type RankedSeason struct {
	Teams         teamMap
	Games         gameMap
	Weeks         WeekList
	WeightedWeeks WeightedWeekMap
}

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
type Rteam struct {
	Id     int
	Week   int
	Rank   int
	Weight int
	Games  []int
	Stats  TrackedStats
}

type teamMap map[int]RankerTeam
type gameMap map[int]RankerGame
type Week struct {
	Week  int
	Games []int
}

type WeekList []Week

type WkTeamMap map[int]*Rteam

type TeamAndRank struct {
	Id     int
	Rank   int
	Weight int
}

type WeightedWeekMap []WkTeamMap

type rankerParams = struct {
	tms      []*Rteam
	sortFn   (func(i, j int) bool)
	accessor func(*Rteam) int
	assigner func(*Rteam, int)
}

/* Fns */
/* Top Level Ranking Fn */
func Rank(
	teams []RankerTeam,
	games []RankerGame,
) (RankedSeason, error) {
	szn := BuildSeason(teams, games)
	err := szn.CompileSeasonStats()
	if err != nil {
		return RankedSeason{}, err
	}
	szn.CalculateStatRankings()
	return szn, nil
}

func BuildSeason(teams []RankerTeam, games []RankerGame) RankedSeason {

	tm := teamMap{}
	for _, t := range teams {
		tm[t.Id] = t
	}

	gmMap := gameMap{}
	wl := WeekList{}

	for _, g := range games {
		gmMap[g.Id] = g

		if g.Week >= len(wl) {
			tmp := make(WeekList, g.Week+1)
			copy(tmp, wl)
			wl = tmp
		}
		if wl[g.Week].Games == nil {

			wl[g.Week] = Week{
				Week:  g.Week,
				Games: []int{},
			}
		}

		wk := wl[g.Week]
		wk.Games = append(wk.Games, g.Id)
		wl[g.Week] = wk
	}

	weightMap := WeightedWeekMap{}

	return RankedSeason{
		Teams:         tm,
		Games:         gmMap,
		Weeks:         wl,
		WeightedWeeks: weightMap,
	}
}

/* Ranking Methods */
func (rs *RankedSeason) CompileSeasonStats() error {

	var UpdateWeightedTeam = func(currTeam RankerStat, oppTeam RankerStat, tm *Rteam) error {
		if currTeam.Id != tm.Id {
			return fmt.Errorf("invalid team ids (%v-%v)", currTeam.Id, tm.Id)
		}

		//Points
		tm.Stats.PF.Val += currTeam.Points
		tm.Stats.PA.Val += oppTeam.Points

		//Yards
		tm.Stats.TotalOffense.Val += currTeam.TotalYards
		tm.Stats.TotalDefense.Val += oppTeam.TotalYards

		//W/L
		if currTeam.Points > oppTeam.Points {
			tm.Stats.Wins.Val++
		} else {
			tm.Stats.Losses.Val++
		}

		return nil
	}

	for _, wk := range rs.Weeks {

		if len(rs.WeightedWeeks) == 0 {
			//Create first week if none exist
			wkMap := make(map[int]*Rteam)
			for _, tm := range rs.Teams {
				newTm := Rteam{
					Id:   tm.Id,
					Week: wk.Week,
				}
				wkMap[tm.Id] = &newTm
			}
			rs.WeightedWeeks = append(rs.WeightedWeeks, wkMap)
		} else {

			//copy prev week into curr week
			cpy, err := rs.WeightedWeeks.copyWeek(wk.Week-1, wk.Week)
			if err != nil {
				return err
			}
			rs.WeightedWeeks = append(rs.WeightedWeeks, cpy)
		}

		for _, gmId := range wk.Games {

			gm, ok := rs.Games[gmId]
			if !ok {
				return fmt.Errorf("game id (%v) not found", gmId)
			}

			homeTeam := gm.Stats.Home
			awayTeam := gm.Stats.Away

			//Home team stats
			if wtHome, ok := rs.WeightedWeeks[wk.Week][homeTeam.Id]; ok {
				err := UpdateWeightedTeam(homeTeam, awayTeam, wtHome)
				if err != nil {
					return err
				}
				wtHome.Games = append(wtHome.Games, gmId)
				rs.WeightedWeeks[wk.Week][homeTeam.Id] = wtHome

			}

			//away team stats
			if wtAway, ok := rs.WeightedWeeks[wk.Week][awayTeam.Id]; ok {
				err := UpdateWeightedTeam(awayTeam, homeTeam, wtAway)
				if err != nil {
					return err
				}
				wtAway.Games = append(wtAway.Games, gmId)
				rs.WeightedWeeks[wk.Week][awayTeam.Id] = wtAway
			}
		}
	}
	return nil
}

func (rs *RankedSeason) CalculateStatRankings() {

	//iterate through weightedWeeks and sort stats assign rank
	for _, wk := range rs.WeightedWeeks {

		tms := wk.toSlice()

		rankConfigs := []rankerParams{
			//wins
			makeRanker(
				tms,
				func(t *Rteam) int { return t.Stats.Wins.Val },
				func(t *Rteam, i int) { t.Stats.Wins.Rank = i },
				true,
			),
			// off
			makeRanker(
				tms,
				func(t *Rteam) int { return t.Stats.TotalOffense.Val },
				func(t *Rteam, i int) { t.Stats.TotalOffense.Rank = i },
				true,
			),
			//pf
			makeRanker(tms,
				func(t *Rteam) int { return t.Stats.PF.Val },
				func(t *Rteam, i int) { t.Stats.PF.Rank = i },
				true,
			),
			/* Descending stats */

			//losses
			makeRanker(
				tms,
				func(t *Rteam) int { return t.Stats.Losses.Val },
				func(t *Rteam, i int) { t.Stats.Losses.Rank = i },
				false,
			),
			//def
			makeRanker(
				tms,
				func(t *Rteam) int { return t.Stats.TotalDefense.Val },
				func(t *Rteam, i int) { t.Stats.TotalDefense.Rank = i },
				false,
			),
			//PA
			makeRanker(
				tms,
				func(t *Rteam) int { return t.Stats.PA.Val },
				func(t *Rteam, i int) { t.Stats.PA.Rank = i },
				false,
			),
		}

		for _, cfg := range rankConfigs {
			rankStat(cfg)
		}

		//sum weights
		for _, tms := range rs.WeightedWeeks {
			for _, tm := range tms {
				tm.sumWeights()
			}
		}

		// Assign overall ranking
		for _, tms := range rs.WeightedWeeks {
			tmSlice := tms.toSlice()

			sort.Slice(tmSlice, func(i, j int) bool {
				return tmSlice[i].Weight < tmSlice[j].Weight
			})

			rankStat(
				makeRanker(
					tmSlice,
					func(t *Rteam) int { return t.Weight },
					func(t *Rteam, i int) { t.Rank = i },
					false,
				))
		}
	}
}

/* Helper Methods */
func (wk *WkTeamMap) toSlice() []*Rteam {
	tms := make([]*Rteam, len(*wk))
	i := 0
	for _, v := range *wk {
		tms[i] = v
		i++
	}
	return tms
}

func (wk *WkTeamMap) GetRankings() ([]TeamAndRank, error) {
	tmPtrs := wk.toSlice()
	tms := make([]TeamAndRank, len(tmPtrs))
	for i, t := range tmPtrs {
		tms[i] = TeamAndRank{
			Id:     t.Id,
			Rank:   t.Rank,
			Weight: t.Weight,
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

	return tms, nil
}

func (wkMap *WeightedWeekMap) copyWeek(week int, newWeek int) (map[int]*Rteam, error) {
	if week < 0 {
		return nil, errors.New("week cannot be less than 0")
	}
	mapCopy := map[int]*Rteam{}

	for k, v := range (*wkMap)[week] {
		gamesCopy := make([]int, len(v.Games))
		copy(gamesCopy, v.Games)

		mapCopy[k] = &Rteam{
			Id:     v.Id,
			Week:   newWeek,
			Rank:   v.Rank,
			Weight: v.Weight,
			Games:  gamesCopy,
			Stats:  v.Stats,
		}
	}

	return mapCopy, nil
}

/* TODO will need to add weight param */
func (tm *Rteam) sumWeights() {
	wt := 0
	wt += tm.Stats.Wins.Rank
	wt += tm.Stats.Losses.Rank
	wt += tm.Stats.TotalOffense.Rank
	wt += tm.Stats.TotalDefense.Rank
	wt += tm.Stats.PF.Rank
	wt += tm.Stats.PA.Rank

	tm.Weight = wt
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
	tms []*Rteam,
	accessor func(*Rteam) int,
	assigner func(*Rteam, int),
	desc bool,
) rankerParams {
	return rankerParams{
		tms:      tms,
		sortFn:   sortTeams(tms, accessor, desc),
		accessor: accessor,
		assigner: assigner,
	}
}

func sortTeams(tms []*Rteam, accessor func(t *Rteam) int, desc bool) func(i, j int) bool {
	return func(i, j int) bool {
		a := accessor(tms[i])
		b := accessor(tms[j])
		if desc {
			return a > b
		}
		return a < b
	}
}
