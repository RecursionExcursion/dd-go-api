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
	buildSeason(teams, games)
}

type team struct {
	id     int
	week   int
	rank   int
	weight int
	games  []int
	stats  struct {
		Wins         Stat
		Losses       Stat
		TotalOffense Stat
		TotalDefense Stat
		PF           Stat
		PA           Stat
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

type season struct {
	teams         teamMap
	games         gameMap
	weeks         weekList
	weightedWeeks weightedWeekMap
}

func buildSeason(teams []RankerTeam, games []RankerGame) {

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

	log.Println(tm)
	log.Println(gm)
	log.Println(wl)
	log.Println(weightMap)

}
