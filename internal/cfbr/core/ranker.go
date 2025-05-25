package core

import "log"

type Ranker struct {
	// rankerSeason RankerSeason
}

type RankerTeam = struct {
	id int
}

type Stat = struct {
	id         int
	totalYards int
	points     int
}

type RankerGameStats = struct {
	home Stat
	away Stat
}

type RankerGame = struct {
	id    int
	week  int
	stats RankerGameStats
}

// type RankerSeason = struct {
// 	// options struct {
// 	// }
// }

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
		tm[t.id] = t
	}

	gm := gameMap{}
	wl := weekList{}
	for _, g := range games {
		gm[g.id] = g

		if g.week > len(wl) {
			tmp := make(weekList, g.week)
			copy(tmp, wl)
			wl = tmp
			wl[g.week] = week{
				week:  g.week,
				games: []int{},
			}
		}
		wk := wl[g.week]
		wk.games = append(wk.games, g.id)
		wl[g.week] = wk
	}

	weightMap := weightedWeekMap{}

	log.Println(tm)
	log.Println(gm)
	log.Println(wl)
	log.Println(weightMap)

}
