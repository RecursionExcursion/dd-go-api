package core

import (
	"fmt"
	"strconv"
)

type SerializeableCompressedSeason struct {
	Id               string `json:"id"`
	CreatedAt        int    `json:"createdAt"`
	Year             int    `json:"year"`
	CompressedSeason string `json:"season"`
}

func MapToRanker(szn Season) (tms []RankerTeam, gms []RankerGame, err error) {

	//Map szsn to ds's
	tms = make([]RankerTeam, len(szn.Teams))
	i := 0
	for _, tm := range szn.Teams {
		id, err := strconv.Atoi(tm.Id)
		if err != nil {
			return tms, gms, fmt.Errorf("could not cast tm id (%v) to int", tm.Id)
		}
		tms[i] = RankerTeam{
			Id: id,
		}
		i++
	}

	gms = make([]RankerGame, len(szn.Games))
	i = 0

	for _, g := range szn.Games {

		//cast id
		id, err := strconv.Atoi(g.Header.Id)
		if err != nil {
			return tms, gms, err
		}

		homeStats := RankerStat{}
		awayStats := RankerStat{}

		for _, tm := range g.Boxscore.Teams {
			if tm.HomeAway == "home" {
				homeStats, err = tmToStat(tm, g)
				if err != nil {
					return tms, gms, err
				}
			} else if tm.HomeAway == "away" {
				awayStats, err = tmToStat(tm, g)
				if err != nil {
					return tms, gms, err
				}
			} else {
				return tms, gms, fmt.Errorf("homeAway not found in game (%v)", g.Header.Id)
			}
		}

		gms[i] = RankerGame{
			Id:   id,
			Week: g.Header.Week,
			Type: g.Header.Season.Type,
			Stats: RankerGameStats{
				Home: homeStats,
				Away: awayStats,
			},
		}
		i++
	}

	return tms, gms, err
}

func tmToStat(tm Team, gm ESPNCfbGame) (rs RankerStat, err error) {

	rs = RankerStat{}

	id, err := strconv.Atoi(tm.Team.Id)
	if err != nil {
		return rs, err
	}
	rs.Id = id

	for _, st := range tm.Statistics {
		if st.Name == totalYardsStatKey {
			ty, err := strconv.Atoi(st.DisplayValue)
			if err != nil {
				return rs, fmt.Errorf("could not cast total yards to int for game (%v)", gm.Header.Id)
			}
			rs.TotalYards = ty
		}
	}

	for _, c := range gm.Header.Competitions[0].Competitors {
		if c.Id == tm.Team.Id {
			score, err := strconv.Atoi(c.Score)
			if err != nil {
				return rs, err
			}
			rs.Points = score
		}
	}

	return rs, err
}
