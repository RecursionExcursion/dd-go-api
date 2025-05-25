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

func TestConstructSeason(t *testing.T) {

	Rank(mockTeams, mockGames)

	t.Error("sdsd")

}
