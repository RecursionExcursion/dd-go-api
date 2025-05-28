package core

import (
	"testing"
)

func createGame(
	id int,
	wk int,
	home RankerStat,
	away RankerStat,
) RankerGame {
	return RankerGame{
		Id:   id,
		Week: wk,
		Stats: RankerGameStats{
			Home: home,
			Away: away,
		},
	}
}

var mockTeams = []RankerTeam{
	{Id: 1},
	{Id: 2},
	{Id: 3},
}

var mockGames = []RankerGame{

	/* Week 1 */
	createGame(11, 1,
		RankerStat{
			Id:         1,
			TotalYards: 100,
			Points:     7,
		},
		RankerStat{
			Id:         2,
			TotalYards: 75,
			Points:     5,
		},
	),

	createGame(
		12,
		1,
		RankerStat{
			Id:         3,
			TotalYards: 150,
			Points:     10,
		},
		RankerStat{
			Id:         4,
			TotalYards: 25,
			Points:     0,
		},
	),

	/* week 2 */

	createGame(
		13,
		2,
		RankerStat{
			Id:         1,
			TotalYards: 100,
			Points:     7,
		},
		RankerStat{
			Id:         3,
			TotalYards: 150,
			Points:     10,
		},
	),

	createGame(
		14,
		2,
		RankerStat{
			Id:         2,
			TotalYards: 75,
			Points:     5,
		},
		RankerStat{
			Id:         4,
			TotalYards: 25,
			Points:     0,
		},
	),
}

func TestBuildSeason(t *testing.T) {

	szn := BuildSeason(mockTeams, mockGames)

	//test all teams were added, t4 is lower div so not checked
	for k, v := range szn.teams {
		if !(k == 1 || k == 2 || k == 3) {
			t.Errorf("Invalid team id (%v)", k)
		}
		if szn.teams[k] != v {
			t.Errorf("Team id (%v) key (%v) mismatch", v, k)
		}
	}

	//test all games were added
	for k := range szn.games {
		if !(k == 11 || k == 12 || k == 13 || k == 14) {
			t.Errorf("Invalid game id (%v)", k)
		}
	}

	//test week list
	wkList := szn.weeks
	if len(wkList) != 3 {
		t.Errorf("Week list is of len %v instead of %v", len(wkList), 3)
	}
	if len(wkList[0].games) > 0 {
		t.Errorf("Week 0 should have 0 entires but instead has %v", wkList[0].games)
	}
	if len(wkList[1].games) != 2 {
		t.Errorf("Week 1 should have 2 entires but instead has %v", wkList[1].games)
	}
	if len(wkList[2].games) != 2 {
		t.Errorf("Week 2 should have 2 entires but instead has %v", wkList[2].games)
	}

	// log.Println(szn)
}

type expectedStats struct {
	wins int
	loss int
	off  int
	def  int
	pa   int
	pf   int
}

func TestCompileSeasonStats(t *testing.T) {
	var checkStats = func(expected expectedStats, tm team) {
		tmId := tm.id
		stats := tm.stats

		if stats.Wins.Val != expected.wins {
			t.Errorf("Team (%v) expected %v wins but had %v", tmId, expected.wins, stats.Wins)
		}

		if stats.Losses.Val != expected.loss {
			t.Errorf("Team (%v) expected %v losses but had %v", tmId, expected.loss, stats.Losses)
		}

		if stats.TotalOffense.Val != expected.off {
			t.Errorf("Team (%v) expected %v off but had %v", tmId, expected.off, stats.TotalOffense)
		}

		if stats.TotalDefense.Val != expected.def {
			t.Errorf("Team (%v) expected %v def but had %v", tmId, expected.def, stats.TotalDefense)
		}

		if stats.PF.Val != expected.pf {
			t.Errorf("Team (%v) expected %v pf but had %v", tmId, expected.pf, stats.PF)
		}

		if stats.PA.Val != expected.pa {
			t.Errorf("Team (%v) expected %v pa but had %v", tmId, expected.pa, stats.PA)
		}
	}

	rs := BuildSeason(mockTeams, mockGames)
	CompileSeasonStats(&rs)

	//check tm1 wk1
	checkStats(expectedStats{
		wins: 1,
		loss: 0,
		off:  100,
		def:  75,
		pf:   7,
		pa:   5,
	},
		rs.weightedWeeks[1][1],
	)

	//check tm1 wk2
	checkStats(expectedStats{
		wins: 1,
		loss: 1,
		off:  200,
		def:  225,
		pf:   14,
		pa:   15,
	},
		rs.weightedWeeks[2][1],
	)

	//check tm3 wk1
	checkStats(expectedStats{
		wins: 1,
		loss: 0,
		off:  150,
		def:  25,
		pf:   10,
		pa:   0,
	},
		rs.weightedWeeks[1][3],
	)

	//check tm3 wk2
	checkStats(expectedStats{
		wins: 2,
		loss: 0,
		off:  300,
		def:  125,
		pf:   20,
		pa:   7,
	},
		rs.weightedWeeks[2][3],
	)
}

func TestCalculateStatRankings(t *testing.T) {
	var checkRank = func(expected int, actual int) {
		if expected != actual {
			t.Errorf("Ranking Error: Expected %v but got %v", expected, actual)
		}
	}

	rs := BuildSeason(mockTeams, mockGames)
	CompileSeasonStats(&rs)
	CalculateStatRankings(&rs)

	t3w2 := rs.weightedWeeks[2][3]

	checkRank(1, t3w2.stats.Wins.Rank)
	checkRank(1, t3w2.stats.TotalOffense.Rank)

	t1w2 := rs.weightedWeeks[2][1]
	checkRank(2, t1w2.stats.TotalOffense.Rank)

}
