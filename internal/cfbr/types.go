package cfbr

/* Retrieval types */

type Team struct {
	Id             uint     `json:"id"`
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
	Id         uint   `json:"id"`
	Season     uint   `json:"season"`
	Week       uint   `json:"week"`
	SeasonType string `json:"season_type"`
	StartDate  string `json:"start_date"`
	Completed  bool   `json:"completed"`
	HomeId     uint   `json:"home_id"`
	HomePoints uint   `json:"home_points"`
	AwayId     uint   `json:"away_id"`
	AwayPoints uint   `json:"away_points"`
}

type GameStats struct {
	Id    uint       `json:"id"`
	Teams []GameTeam `json:"teams"`
}

type GameTeam struct {
	SchoolId uint   `json:"schoolId"`
	HomeAway string `json:"homeAway"`
	Points   uint   `json:"points"`
	Stats    []struct {
		Category string `json:"category"`
		Stat     string `json:"stat"`
	} `json:"stats"`
}
