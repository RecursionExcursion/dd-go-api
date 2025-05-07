package core

// import (
// 	"strconv"
// 	"testing"
// )

// func TestCompileGameStats(t *testing.T) {

// 	/* Mock Values */
// 	off := 100
// 	def := 50
// 	pf := 10
// 	pa := 7
// 	w := 1
// 	l := 0

// 	mockCompleteGame := CompleteGame{
// 		Id:   69,
// 		Game: Game{},
// 		GameStats: GameStats{
// 			Id: 69,
// 			Teams: []GameTeam{
// 				{
// 					SchoolId: 1,
// 					HomeAway: "home",
// 					Points:   pf,
// 					Stats: []GameTeamStats{
// 						{
// 							Category: totalYardsStatKey,
// 							Stat:     strconv.Itoa(off),
// 						},
// 					},
// 				},
// 				{
// 					SchoolId: 2,
// 					HomeAway: "away",
// 					Points:   pa,
// 					Stats: []GameTeamStats{
// 						{
// 							Category: totalYardsStatKey,
// 							Stat:     strconv.Itoa(def),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	ts, err := compileGameStats(1, mockCompleteGame)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if ts.Total.Wins.Value != w {
// 		t.Errorf("Expected Wins to be %v but got %v", w, ts.Total.Wins.Value)
// 	}
// 	if ts.Total.Losses.Value != l {
// 		t.Errorf("Expected Losses to be %v but got %v", l, ts.Total.Losses.Value)
// 	}

// 	if ts.Total.TotalOffense.Value != off {
// 		t.Errorf("Expected Wins to be %v but got %v", off, ts.Total.TotalOffense.Value)
// 	}
// 	if ts.Total.TotalDefense.Value != def {
// 		t.Errorf("Expected Wins to be %v but got %v", def, ts.Total.TotalDefense.Value)
// 	}

// 	if ts.Total.PF.Value != pf {
// 		t.Errorf("Expected Wins to be %v but got %v", pf, ts.Total.PF.Value)
// 	}
// 	if ts.Total.PA.Value != pa {
// 		t.Errorf("Expected Wins to be %v but got %v", pa, ts.Total.PA.Value)
// 	}
// }
