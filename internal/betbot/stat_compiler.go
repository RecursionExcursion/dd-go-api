package betbot

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

type PackagedPlayer struct {
	Name                string
	Team                string
	FirstToScore        uint8
	FirstShotAttempts   uint8
	ScoreOnFirstAttempt uint8
}

type StatCalculator struct {
	fsd FirstShotData
}

func NewStatCalculator(fsd FirstShotData) *StatCalculator {

	filteredGames := []game{}

	for _, g := range fsd.Games {
		if g.Season.Slug == "preseason" {
			continue
		}
		filteredGames = append(filteredGames, g)
	}

	fsd.Games = filteredGames

	FindGameInFsd(fsd, strconv.Itoa(401705613))

	return &StatCalculator{
		fsd,
	}
}

func (sc *StatCalculator) CalcAndPackage() ([]PackagedPlayer, error) {

	log.Printf("Calc stats for %v games", len(sc.fsd.Games))

	err := sc.calcFirstScore()
	if err != nil {
		return []PackagedPlayer{}, err
	}

	err = sc.calcFirstShotAttempt()
	if err != nil {
		return []PackagedPlayer{}, err
	}

	allData := sc.packageData()
	var filteredData = []PackagedPlayer{}

	for _, pp := range allData {
		if pp.ScoreOnFirstAttempt == 0 &&
			pp.FirstShotAttempts == 0 &&
			pp.FirstToScore == 0 {
			continue
		}
		filteredData = append(filteredData, pp)
	}

	sort.Slice(filteredData, func(a, b int) bool {
		return filteredData[a].FirstToScore > filteredData[b].FirstToScore
	})

	return filteredData, nil
}

func (sc *StatCalculator) calcFirstScore() error {
	for _, gm := range sc.fsd.Games {

		fs := gm.TrackedEvents.FirstScore

		if fs.Id == "" {
			err := fmt.Errorf("first shot data for game:%v not found", gm.Id)
			lib.LogError("", err)
			continue
		}

		playerId := fs.Participants[0].Athlete.Id

		player, err := sc.findPlayerById(playerId)
		if err != nil {
			lib.LogError(fmt.Sprintf("Player %v for play %v not found", playerId, fs.Text), err)
			continue
		}
		player.BetStats.FirstPointsMade++
	}
	return nil
}

func (sc *StatCalculator) calcFirstShotAttempt() error {
	for _, gm := range sc.fsd.Games {

		fsa := gm.TrackedEvents.FirstShotAttempt

		if fsa.Id == "" {
			err := fmt.Errorf("first shot attempt data for game:%v not found", gm.Id)
			lib.LogError("", err)
			continue
		}

		playerId := fsa.Participants[0].Athlete.Id
		player, err := sc.findPlayerById(playerId)
		if err != nil {
			lib.LogError(fmt.Sprintf("Player %v for play %v not found", playerId, fsa.Text), err)
			continue
		}
		player.BetStats.FirstShotAttempts++
		if fsa.ScoringPlay {
			player.BetStats.ScoreOnFirstAttempt++
		}
	}
	return nil
}

func (sc *StatCalculator) packageData() []PackagedPlayer {
	var packagedPlayers = []PackagedPlayer{}

	for _, t := range sc.fsd.Teams {
		for _, p := range t.Roster {

			pp := PackagedPlayer{
				Name:                p.FullName,
				Team:                t.Name,
				FirstToScore:        p.BetStats.FirstPointsMade,
				FirstShotAttempts:   p.BetStats.FirstShotAttempts,
				ScoreOnFirstAttempt: p.BetStats.ScoreOnFirstAttempt,
			}

			packagedPlayers = append(packagedPlayers, pp)
		}
	}
	return packagedPlayers
}

func (sc *StatCalculator) findPlayerById(id string) (*player, error) {
	var player player

	for _, t := range sc.fsd.Teams {
		for i := range t.Roster {
			if t.Roster[i].Id == id {
				return &t.Roster[i], nil
			}
		}
	}
	return &player, fmt.Errorf("player %v not found", id)
}
