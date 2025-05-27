package core

import (
	"log"
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

func TestRankSeason(t *testing.T) {
	szn := BuildSeason(mockTeams, mockGames)
	rs := RankSeason(&szn)
	log.Println(rs)

}
