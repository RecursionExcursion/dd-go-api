package pickle

import "slices"

type PicklePlayer struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Matches []string `json:"matches"`
}

func (pp *PicklePlayer) removeMatch(id string) bool {
	for i, mId := range pp.Matches {
		if mId == id {
			pp.Matches = slices.Delete(pp.Matches, i, i+1)
			return true
		}
	}
	return false
}

type PickleMatch struct {
	Id    string       `json:"id"`
	Date  int          `json:"date"`
	Score []MatchScore `json:"score"`
}

type MatchScore struct {
	Id     string `json:"id"`
	Points int    `json:"points"`
}

type PickleData struct {
	ID      string         `json:"id"`
	Players []PicklePlayer `json:"players"`
	Matches []PickleMatch  `json:"matches"`
}

func (pd *PickleData) addPlayer(np PicklePlayer) {
	pd.Players = append(pd.Players, np)
}

func (pd *PickleData) getPlayer(id string) (*PicklePlayer, bool) {
	var nilPlayer PicklePlayer

	for i := range pd.Players {
		if pd.Players[i].Id == id {
			return &pd.Players[i], true
		}
	}
	return &nilPlayer, false
}

func (pd *PickleData) removePlayer(id string) bool {
	for i, p := range pd.Players {
		if p.Id == id {
			pd.Players = slices.Delete(pd.Players, i, i+1)
			return true
		}
	}
	return false
}

func (pd *PickleData) addMatch(pm PickleMatch) bool {
	//trim pto 2 players
	pm.Score = []MatchScore{
		pm.Score[0],
		pm.Score[1],
	}

	p1Id := pm.Score[0].Id
	p2Id := pm.Score[1].Id

	p1, ok := pd.getPlayer(p1Id)
	if !ok {
		return false
	}

	p2, ok := pd.getPlayer(p2Id)
	if !ok {
		return false
	}

	p1.Matches = append(p1.Matches, pm.Id)
	p2.Matches = append(p2.Matches, pm.Id)
	pd.Matches = append(pd.Matches, pm)
	return true
}

func (pd *PickleData) removeMatch(id string) bool {
	for i, m := range pd.Matches {
		if m.Id == id {
			p1, _ := pd.getPlayer(m.Score[0].Id)
			p2, _ := pd.getPlayer(m.Score[1].Id)
			p1.removeMatch(m.Id)
			p2.removeMatch(m.Id)
			pd.Matches = slices.Delete(pd.Matches, i, i+1)
			return true
		}
	}
	return false
}
