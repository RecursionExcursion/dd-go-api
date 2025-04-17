package cfbr

import "log"

type TrackedStats struct {
	Total struct {
		TotalOffense int
		TotalDefense int
		PF           int
		PA           int
	}
	PG struct {
		OffPG int
		DefPG int
		PFPG  int
		PAPG  int
	}
}

type WeeklyTeam struct {
	Id          uint
	GamesPlayed []CompleteGame
	Stats       TrackedStats
}

type ComputedSeason struct {
	SeasonInfo         CFBRSeason
	RegularSeasonWeeks [][]WeeklyTeam
	PostSeasonWeeks    [][]WeeklyTeam
}

func ComputeSeason(s CFBRSeason) ComputedSeason {

	cs := ComputedSeason{
		SeasonInfo:         s,
		RegularSeasonWeeks: make([][]WeeklyTeam, FindLastWeek(s)),
		PostSeasonWeeks:    make([][]WeeklyTeam, 1),
	}

	//DO ALOT OF WORK

	return cs
}

func FindLastWeek(s CFBRSeason) int {

	lastRegSeasonWeek := 0

	for _, sch := range s.Schools {
		for _, g := range sch.Games {
			if g.Game.SeasonType == regularSeason &&
				g.Game.Completed &&
				g.Game.Week > uint(lastRegSeasonWeek) {
				lastRegSeasonWeek = int(g.Game.Week)

			}

		}
	}

	log.Println("Last week")

	return lastRegSeasonWeek
}
