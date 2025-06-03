package core

import (
	"errors"
	"fmt"
	"log"
	"slices"
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
	Type  int
}

/* Types for internal calcs */
type WeightedStat struct {
	Rank int
	Val  int
}

// Consider stat map for dynamic stats
// Could have options "damage opposition"
type TrackedStats struct {
	Wins         WeightedStat
	Losses       WeightedStat
	TotalOffense WeightedStat
	TotalDefense WeightedStat
	PF           WeightedStat
	PA           WeightedStat
}

// TODO dont forget PI and SS
type RankableTeam struct {
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

type WkTeamMap map[int]*RankableTeam

type TeamAndRank struct {
	Id     int
	Rank   int
	Weight int
}

type WeightedWeekMap []WkTeamMap

type rankerParams = struct {
	tms      []*RankableTeam
	sortFn   (func(i, j int) bool)
	accessor func(*RankableTeam) int
	assigner func(*RankableTeam, int)
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

	for _, wk := range rs.Weeks {

		if len(rs.WeightedWeeks) == 0 {
			//Create first week if none exist
			wkMap := make(map[int]*RankableTeam)
			for _, tm := range rs.Teams {
				newTm := RankableTeam{
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
				err := UpdateWeightedTeam(wtHome, gm)
				if err != nil {
					return err
				}
				wtHome.Games = append(wtHome.Games, gmId)
				rs.WeightedWeeks[wk.Week][homeTeam.Id] = wtHome

			}

			//away team stats
			if wtAway, ok := rs.WeightedWeeks[wk.Week][awayTeam.Id]; ok {
				err := UpdateWeightedTeam(wtAway, gm)
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
				func(t *RankableTeam) int { return t.Stats.Wins.Val },
				func(t *RankableTeam, i int) { t.Stats.Wins.Rank = i },
				true,
			),
			// off
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.TotalOffense.Val },
				func(t *RankableTeam, i int) { t.Stats.TotalOffense.Rank = i },
				true,
			),
			//pf
			makeRanker(tms,
				func(t *RankableTeam) int { return t.Stats.PF.Val },
				func(t *RankableTeam, i int) { t.Stats.PF.Rank = i },
				true,
			),
			/* Descending stats */

			//losses
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.Losses.Val },
				func(t *RankableTeam, i int) { t.Stats.Losses.Rank = i },
				false,
			),
			//def
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.TotalDefense.Val },
				func(t *RankableTeam, i int) { t.Stats.TotalDefense.Rank = i },
				false,
			),
			//PA
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.PA.Val },
				func(t *RankableTeam, i int) { t.Stats.PA.Rank = i },
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
					func(t *RankableTeam) int { return t.Weight },
					func(t *RankableTeam, i int) { t.Rank = i },
					false,
				))
		}
	}
}

/* Helper Methods */
func (wk *WkTeamMap) toSlice() []*RankableTeam {
	tms := make([]*RankableTeam, len(*wk))
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

func (wkMap *WeightedWeekMap) copyWeek(week int, newWeek int) (map[int]*RankableTeam, error) {
	if week < 0 {
		return nil, errors.New("week cannot be less than 0")
	}
	mapCopy := map[int]*RankableTeam{}

	for k, v := range (*wkMap)[week] {
		gamesCopy := make([]int, len(v.Games))
		copy(gamesCopy, v.Games)

		mapCopy[k] = &RankableTeam{
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
func (tm *RankableTeam) sumWeights() {
	wt := 0
	wt += tm.Stats.Wins.Rank
	wt += tm.Stats.Losses.Rank
	wt += tm.Stats.TotalOffense.Rank
	wt += tm.Stats.TotalDefense.Rank
	wt += tm.Stats.PF.Rank
	wt += tm.Stats.PA.Rank

	tm.Weight = wt
}

func (rt *RankerTeam) teamToRankable() RankableTeam {
	return RankableTeam{
		Id:    rt.Id,
		Games: []int{},
	}
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
	tms []*RankableTeam,
	accessor func(*RankableTeam) int,
	assigner func(*RankableTeam, int),
	desc bool,
) rankerParams {
	return rankerParams{
		tms:      tms,
		sortFn:   sortTeams(tms, accessor, desc),
		accessor: accessor,
		assigner: assigner,
	}
}

func sortTeams(tms []*RankableTeam, accessor func(t *RankableTeam) int, desc bool) func(i, j int) bool {
	return func(i, j int) bool {
		a := accessor(tms[i])
		b := accessor(tms[j])
		if desc {
			return a > b
		}
		return a < b
	}
}

////////////////////////////////////////////////////////////////////////

type SeasonMap map[int]map[int]*RankableTeam

func RankSeasonProto(tms []RankerTeam, gms []RankerGame) SeasonMap {
	//create collections
	tmMap := makeTeamMap(tms)
	// sznMap := make(SeasonMap)
	sznTypeMap := make(map[int]SeasonMap)

	//iterate games and add stats to sznMap
	for _, g := range gms {

		tp := g.Type
		//season type
		sznMap, ok := sznTypeMap[tp]
		if !ok {
			sznTypeMap[tp] = make(SeasonMap)
			sznMap = sznTypeMap[tp]
		}

		//week
		wk := g.Week

		//check map
		wkMap, ok := sznMap[wk]
		if !ok {
			sznMap[wk] = makeRankableTeamMap(tmMap, wk)
			wkMap = sznMap[wk]
		}

		if ht, ok := wkMap[g.Stats.Home.Id]; ok {
			UpdateWeightedTeam(ht, g)
		}
		if at, ok := wkMap[g.Stats.Away.Id]; ok {
			UpdateWeightedTeam(at, g)
		}
	}
	sznMap := squashSeasonMaps(sznTypeMap)

	//squash stats
	squashStats(sznMap)
	calcWeights(sznMap)

	return sznMap
}

func makeTeamMap(tms []RankerTeam) teamMap {
	tmMap := make(teamMap)
	for _, tm := range tms {
		tmMap[tm.Id] = tm
	}
	return tmMap
}

func makeRankableTeamMap(tmMap teamMap, wk int) map[int]*RankableTeam {

	weekMap := make(map[int]*RankableTeam)

	for id, tm := range tmMap {
		val := tm.teamToRankable()
		val.Week = wk
		weekMap[id] = &val
	}

	return weekMap
}

func UpdateWeightedTeam(tm *RankableTeam, gm RankerGame) error {
	tmId := tm.Id

	currTeam := RankerStat{}
	oppTeam := RankerStat{}

	if gm.Stats.Home.Id == tmId {
		currTeam = gm.Stats.Home
		oppTeam = gm.Stats.Away
	} else if gm.Stats.Away.Id == tmId {
		currTeam = gm.Stats.Away
		oppTeam = gm.Stats.Home
	} else {
		return fmt.Errorf("Team %v not found in game %v", tmId, gm.Id)
	}

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

	tm.Games = append(tm.Games, gm.Id)

	return nil
}

func squashStats(sznMap SeasonMap) {

	keys := []int{}

	for k := range sznMap {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	prev := -1
	for _, k := range keys {
		//grab week
		wk := sznMap[k]

		//check if its the first week, if so skip and iterate
		if prev != -1 {
			prevWk := sznMap[prev]

			for _, tm := range wk {
				prevWkTm := prevWk[tm.Id]

				//append stats
				tm.Stats.Wins.Val += prevWkTm.Stats.Wins.Val
				tm.Stats.Losses.Val += prevWkTm.Stats.Losses.Val

				tm.Stats.TotalOffense.Val += prevWkTm.Stats.TotalOffense.Val
				tm.Stats.TotalDefense.Val += prevWkTm.Stats.TotalDefense.Val

				tm.Stats.PF.Val += prevWkTm.Stats.PF.Val
				tm.Stats.PA.Val += prevWkTm.Stats.PA.Val

				tm.Games = append(prevWkTm.Games, tm.Games...)
			}
		}
		prev = k
	}
}

func calcWeights(sznMap SeasonMap) {

	//iterate through weightedWeeks and sort stats assign rank
	for _, wk := range sznMap {

		tms := wkToSlice(wk)

		rankConfigs := []rankerParams{
			//wins
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.Wins.Val },
				func(t *RankableTeam, i int) { t.Stats.Wins.Rank = i },
				true,
			),
			// off
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.TotalOffense.Val },
				func(t *RankableTeam, i int) { t.Stats.TotalOffense.Rank = i },
				true,
			),
			//pf
			makeRanker(tms,
				func(t *RankableTeam) int { return t.Stats.PF.Val },
				func(t *RankableTeam, i int) { t.Stats.PF.Rank = i },
				true,
			),
			/* Descending stats */

			//losses
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.Losses.Val },
				func(t *RankableTeam, i int) { t.Stats.Losses.Rank = i },
				false,
			),
			//def
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.TotalDefense.Val },
				func(t *RankableTeam, i int) { t.Stats.TotalDefense.Rank = i },
				false,
			),
			//PA
			makeRanker(
				tms,
				func(t *RankableTeam) int { return t.Stats.PA.Val },
				func(t *RankableTeam, i int) { t.Stats.PA.Rank = i },
				false,
			),
		}

		for _, cfg := range rankConfigs {
			rankStat(cfg)
		}

		//sum weights
		for _, tms := range sznMap {
			for _, tm := range tms {
				tm.sumWeights()
			}
		}

		// Assign overall ranking
		for _, tms := range sznMap {
			tmSlice := wkToSlice(tms)

			sort.Slice(tmSlice, func(i, j int) bool {
				return tmSlice[i].Weight < tmSlice[j].Weight
			})

			rankStat(
				makeRanker(
					tmSlice,
					func(t *RankableTeam) int { return t.Weight },
					func(t *RankableTeam, i int) { t.Rank = i },
					false,
				))
		}
	}
}

func wkToSlice(wk map[int]*RankableTeam) []*RankableTeam {
	tms := []*RankableTeam{}
	for _, v := range wk {
		tms = append(tms, v)
	}
	return tms
}

func squashSeasonMaps(sznMaps map[int]SeasonMap) SeasonMap {
	types := []int{}

	weekIndex := 0
	finalSzn := make(SeasonMap)

	for k := range sznMaps {
		types = append(types, k)
	}

	slices.Sort(types)

	for _, t := range types {
		currSzn := sznMaps[t]

		wks := []int{}

		for k := range currSzn {
			wks = append(wks, k)
		}

		slices.Sort(wks)

		for _, wk := range wks {
			finalSzn[weekIndex] = currSzn[wk]
			weekIndex++
		}

	}

	log.Println(len(sznMaps))
	log.Println(len(finalSzn))

	return finalSzn
}

/*
ignore 0 week
every week after, eval prev week ranks, rank poll inertia then re weight
*/

func caclPollInertia() {

}

/*
ignore week 0,
eval prev week to calc curr week, do this last
*/
func calcStrengthOfSched() {

}
