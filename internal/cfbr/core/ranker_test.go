package core

import "testing"

var mockTeams = []RankerTeam{
	{id: 1},
	{id: 2},
	{id: 3},
	{id: 4},
}

var mockGames = []RankerGame{
	/* Week 1 */
	{id: 11,
		week: 1,
		stats: RankerGameStats{
			home: Stat{
				id:         1,
				totalYards: 100,
				points:     7,
			},
			away: Stat{
				id:         2,
				totalYards: 75,
				points:     5,
			},
		},
	},

	{id: 12,
		week: 1,
		stats: RankerGameStats{
			home: Stat{
				id:         3,
				totalYards: 150,
				points:     10,
			},
			away: Stat{
				id:         4,
				totalYards: 25,
				points:     0,
			},
		},
	},

	/* week 2 */

	{id: 13,
		week: 2,
		stats: RankerGameStats{
			home: Stat{
				id:         1,
				totalYards: 100,
				points:     7,
			},
			away: Stat{
				id:         3,
				totalYards: 150,
				points:     10,
			},
		},
	},

	{id: 14,
		week: 2,
		stats: RankerGameStats{
			home: Stat{
				id:         2,
				totalYards: 75,
				points:     5,
			},
			away: Stat{
				id:         4,
				totalYards: 25,
				points:     0,
			},
		},
	},
}

func TestConstructSeason(t *testing.T) {

	// ranker := Ranker{struct{

	// 	teams: mockTeams,
	// 	games: mockGames,
	// },
	// }

	Rank(mockTeams, mockGames)

}
