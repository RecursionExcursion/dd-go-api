package core

/* Data Retrieval types */

type Team struct {
	Id             int      `json:"id"`
	School         string   `json:"school"`
	Mascot         string   `json:"mascot"`
	Abbreviation   string   `json:"abbreviation"`
	Conference     string   `json:"conference"`
	Classification string   `json:"classification"`
	Color          string   `json:"color"`
	AltColor       string   `json:"alt_color"`
	Logos          []string `json:"logos"`
}

type Game struct {
	Id         int    `json:"id"`
	Season     int    `json:"season"`
	Week       int    `json:"week"`
	SeasonType string `json:"season_type"`
	StartDate  string `json:"start_date"`
	Completed  bool   `json:"completed"`
	HomeId     int    `json:"home_id"`
	HomePoints int    `json:"home_points"`
	AwayId     int    `json:"away_id"`
	AwayPoints int    `json:"away_points"`
}

type GameStats struct {
	Id    int        `json:"id"`
	Teams []GameTeam `json:"teams"`
}

type GameTeam struct {
	SchoolId int             `json:"schoolId"`
	HomeAway string          `json:"homeAway"`
	Points   int             `json:"points"`
	Stats    []GameTeamStats `json:"stats"`
}

type GameTeamStats struct {
	Category string `json:"category"`
	Stat     string `json:"stat"`
}
