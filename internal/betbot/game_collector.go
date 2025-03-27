package betbot

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func collectGames() ([]game, error) {
	now := time.Now()

	var year int
	if int(now.Month()) >= 8 {
		year = now.Year() + 1
	} else {
		year = now.Year()
	}

	games, err := collectSeasonGames(year)
	if err != nil {
		return nil, err
	}

	fetchPlaysAsync(&games)

	return games, nil
}

func collectSeasonGames(year int) ([]game, error) {

	ranges, err := fetchSeasonInfo(year)
	if err != nil {
		return nil, err
	}

	games, err := fetchSeasonGamesAsync(ranges.Start, ranges.End)
	if err != nil {
		return nil, err
	}

	return games, nil
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func fetchSeasonInfo(year int) (tr TimeRange, err error) {
	yearStr := strconv.Itoa(year)
	seasonEp := endpoints().Season(yearStr)

	fetchFn := func() (*http.Response, error) {
		return http.Get(seasonEp)
	}

	sznPayload, _, err := lib.FetchAndMap[seasonInfoPayload](fetchFn)
	if err != nil {
		return tr, err
	}

	sznInfo := sznPayload.Leagues[0].SeasonInfo

	toTime := func(str string) (time.Time, error) {
		layout := "2006-01-02T15:04Z"
		return time.Parse(layout, str)
	}

	startTime, err := toTime(sznInfo.StartDate)
	if err != nil {
		err = fmt.Errorf("could not parse '%s' to time", sznInfo.StartDate)
		return
	}

	endTime, err := toTime(sznInfo.EndDate)
	if err != nil {
		err = fmt.Errorf("could not parse '%s' to time", sznInfo.EndDate)
		return
	}

	now := time.Now()

	if endTime.After(now) {
		oneDayAgo := now.AddDate(0, 0, -1)
		endTime = oneDayAgo
	}

	tr = TimeRange{
		startTime,
		endTime,
	}

	return tr, nil
}

func fetchSeasonGamesAsync(start time.Time, end time.Time) ([]game, error) {

	var curr = start

	eps := []string{}

	for end.After(curr) {
		fDate := dateToYYYYMMDD(curr)
		sbEp := endpoints().Scoreboard(fDate)
		eps = append(eps, sbEp)
		curr = curr.AddDate(0, 0, 1)
	}

	gChan := make(chan []game, len(eps))
	wg := sync.WaitGroup{}

	gameWorker := func(ep string) {
		defer wg.Done()

		gamesPlayload, _, err := lib.FetchAndMap[seasonGamesFetchPayload](
			func() (*http.Response, error) {
				return http.Get(ep)
			})

		if err != nil {
			gChan <- []game{}
		}

		games := []game{}
		for _, wrapper := range gamesPlayload.Events {
			g := wrapper.game
			//Extract playByPlay bool from nested obj
			if len(wrapper.Competitions) > 0 {
				g.PlayByPlay = wrapper.Competitions[0].PlayByPlayAvailable
			}
			games = append(games, g)
		}

		gChan <- games
	}

	for _, ep := range eps {
		wg.Add(1)
		go gameWorker(ep)
	}

	go func() {
		wg.Wait()
		close(gChan)
	}()

	mappedGames := [][]game{}
	for games := range gChan {
		mappedGames = append(mappedGames, games)
	}

	return flatten2DSlice(mappedGames), nil

}

func dateToYYYYMMDD(d time.Time) string {
	return d.Format("20060102")
}

func fetchPlaysAsync(games *[]game) error {

	type GameDataChannelDTO struct {
		gameData gameDataFetchPayload
		gameId   string
	}

	gdChan := make(chan GameDataChannelDTO)
	wg := sync.WaitGroup{}

	gameDataWorker := func(gameId string) {
		defer wg.Done()

		gEp := endpoints().GameData(gameId)

		gData, _, err := lib.FetchAndMap[gameDataFetchPayload](func() (*http.Response, error) {
			return http.Get(gEp)
		})
		if err != nil {
			log.Println(err)

			gdChan <- GameDataChannelDTO{
				gameId:   gameId,
				gameData: gameDataFetchPayload{},
			}
			return
		}
		gdChan <- GameDataChannelDTO{
			gameId:   gameId,
			gameData: gData,
		}
	}

	//spawn goroutines
	for _, g := range *games {

		if g.Season.Slug == "preseason" || !g.PlayByPlay {
			continue
		}

		gId := g.Id
		wg.Add(1)
		go gameDataWorker(gId)
	}

	go func() {
		wg.Wait()
		close(gdChan)
	}()

	for gdDTO := range gdChan {
		for i := range *games {
			g := &(*games)[i]
			if g.Id == gdDTO.gameId {

				firstPointsPlay, err := extractFirstPoints(gdDTO.gameData.Plays)
				if err != nil {
					log.Printf("%v for game:%v", err.Error(), gdDTO.gameId)
				}

				g.TrackedEvents.FirstScore = firstPointsPlay
				firstShotPlay, err := extractFirstShotAttempt(gdDTO.gameData.Plays)
				if err != nil {
					log.Printf("%v for game:%v", err.Error(), gdDTO.gameId)
				}
				g.TrackedEvents.FirstShotAttempt = firstShotPlay

				break
			}
		}
	}

	return nil
}

func flatten2DSlice[T any](s [][]T) []T {
	res := []T{}
	for _, item := range s {
		res = append(res, item...)
	}
	return res
}
