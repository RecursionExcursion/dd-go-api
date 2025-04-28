package core

import (
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
	mu := sync.Mutex{}
	tasks := []func(){}

	for i := range *teams {
		t := &(*teams)[i]
		teamId := t.Id

		tasks = append(tasks, func() {

			rosterEp := endpoints().Roster(teamId)

			rosterPayload, _, err := lib.FetchAndMap[rosterFetchPayload](
				func() (*http.Response, error) {
					return http.Get(rosterEp)
				})
			if err != nil {
				rChan <- RosterChannel{
					teamId: teamId,
					roster: []player{},
					err:    err,
				}
				return
			}

			rChan <- RosterChannel{
				teamId: teamId,
				roster: rosterPayload.Athletes,
				err:    nil,
			}
		})
	}

	go func() {
		// lib.RunBatch(tasks, batchSize)
		BatchRunner(tasks)
		close(rChan)
	}()

	//ranging a channel is a blocking operation until channel is empty
	for res := range rChan {
		if res.err != nil {
			lib.LogError("error while fetching rosters", res.err)
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
