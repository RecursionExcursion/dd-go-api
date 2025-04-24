package cfbr

import "log"

type TrackedStats struct {
	Total struct {
		Wins         int
		Losses       int
		TotalOffense int
		TotalDefense int
		PF           int
		PA           int
	}
	PG struct {
		WinsPG   int
		LossesPG int
		OffPG    int
		DefPG    int
		PFPG     int
		PAPG     int
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

	for _, g := range s.Games {
		if g.Game.SeasonType == regularSeason &&
			g.Game.Completed &&
			g.Game.Week > uint(lastRegSeasonWeek) {
			lastRegSeasonWeek = int(g.Game.Week)
		}
	}

	log.Printf("Last Reg szn week %v", lastRegSeasonWeek)

	return lastRegSeasonWeek
}
