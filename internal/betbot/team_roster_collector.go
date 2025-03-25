package betbot

import (
	"log"
	"net/http"
	"sync"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func collectTeamsAndRosters() ([]team, error) {

	teams, err := fetchTeams()
	if err != nil {
		return nil, err
	}

	compileErr := compileRosterAsync(&teams)
	if compileErr != nil {
		return nil, compileErr
	}

	return teams, nil
}

func fetchTeams() ([]team, error) {
	tfp, _, err := lib.FetchAndMap[teamFetchPayload](
		func() (*http.Response, error) {
			return http.Get(endpoints().Teams())
		})
	if err != nil {
		return nil, err
	}
	teams := []team{}

	for _, t := range tfp.Sports[0].Leagues[0].Teams {
		teams = append(teams, t.Team)
	}

	return teams, nil
}

func compileRosterAsync(teams *[]team) error {
	type RosterChannel struct {
		teamId string
		roster []player
		err    error
	}

	rChan := make(chan RosterChannel, len(*teams))
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	//TODO consider making inline
	rosterWorker := func(reschan chan RosterChannel,
		teamId string,
	) {
		defer wg.Done()

		rosterEp := endpoints().Roster(teamId)

		rosterPayload, _, err := lib.FetchAndMap[rosterFetchPayload](
			func() (*http.Response, error) {
				return http.Get(rosterEp)
			})
		if err != nil {
			reschan <- RosterChannel{
				teamId: teamId,
				roster: []player{},
				err:    err,
			}
			return
		}

		reschan <- RosterChannel{
			teamId: teamId,
			roster: rosterPayload.Athletes,
			err:    nil,
		}

	}

	for i := range *teams {
		t := &(*teams)[i]
		wg.Add(1)
		go rosterWorker(rChan, t.Id)
	}

	go func() {
		wg.Wait()
		close(rChan)
	}()

	//ranging a channel is a blocking operation until channel is empty
	for res := range rChan {
		if res.err != nil {
			log.Println("error while fetching rosters")
			continue
		}

		mu.Lock()
		for i := range *teams {
			t := &(*teams)[i]
			if t.Id == res.teamId {
				t.Roster = res.roster
				break
			}
		}
		mu.Unlock()

	}

	return nil
}
