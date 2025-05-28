package core

import (
	"testing"
)

func createGame(
	id int,
	wk int,
	home Stat,
	away Stat,
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
	{Id: 4},
}

var mockGames = []RankerGame{

	/* Week 1 */
	createGame(11, 1,
		Stat{
			Id:         1,
			TotalYards: 100,
			Points:     7,
		},
		Stat{
			Id:         2,
			TotalYards: 75,
			Points:     5,
		},
	),

	createGame(
		12,
		1,
		Stat{
			Id:         3,
			TotalYards: 150,
			Points:     10,
		},
		Stat{
			Id:         4,
			TotalYards: 25,
			Points:     0,
		},
	),

	/* week 2 */

	createGame(
		13,
		2,
		Stat{
			Id:         1,
			TotalYards: 100,
			Points:     7,
		},
		Stat{
			Id:         3,
			TotalYards: 150,
			Points:     10,
		},
	),

	createGame(
		14,
		2,
		Stat{
			Id:         2,
			TotalYards: 75,
			Points:     5,
		},
		Stat{
			Id:         4,
			TotalYards: 25,
			Points:     0,
		},
	),
}

func TestBuildSeason(t *testing.T) {

	szn := BuildSeason(mockTeams, mockGames)

	//test all teams were added
	for k, v := range szn.teams {
		if !(k == 1 || k == 2 || k == 3 || k == 4) {
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

func TestCompileSeasonStats(t *testing.T) {
	szn := BuildSeason(mockTeams, mockGames)
	rs := CompileSeasonStats(&szn)

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
		t,
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
		t,
	)

	//check tm4 wk1
	checkStats(expectedStats{
		wins: 1,
		loss: 0,
		off:  150,
		def:  25,
		pf:   10,
		pa:   0,
	},
		rs.weightedWeeks[1][3],
		t,
	)

	//check tm4 wk2
	checkStats(expectedStats{
		wins: 2,
		loss: 0,
		off:  300,
		def:  125,
		pf:   20,
		pa:   7,
	},
		rs.weightedWeeks[2][3],
		t,
	)
}

type expectedStats struct {
	wins int
	loss int
	off  int
	def  int
	pa   int
	pf   int
}

func checkStats(expected expectedStats, tm team, t *testing.T) {
	tmId := tm.id
	stats := tm.stats

	if stats.Wins != expected.wins {
		t.Errorf("Team (%v) expected %v wins but had %v", tmId, expected.wins, stats.Wins)
	}

	if stats.Losses != expected.loss {
		t.Errorf("Team (%v) expected %v losses but had %v", tmId, expected.loss, stats.Losses)
	}

	if stats.TotalOffense != expected.off {
		t.Errorf("Team (%v) expected %v off but had %v", tmId, expected.off, stats.TotalOffense)
	}

	if stats.TotalDefense != expected.def {
		t.Errorf("Team (%v) expected %v def but had %v", tmId, expected.def, stats.TotalDefense)
	}

	if stats.PF != expected.pf {
		t.Errorf("Team (%v) expected %v pf but had %v", tmId, expected.pf, stats.PF)
	}

	if stats.PA != expected.pa {
		t.Errorf("Team (%v) expected %v pa but had %v", tmId, expected.pa, stats.PA)
	}

}
