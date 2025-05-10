package core

/* Data Retrieval types (CFBAPI)*/

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

// ESPN TYPES

/* Groups */

type ESPNGroups struct {
	Status string  `json:"status"`
	Groups []Group `json:"groups"`
}

type Group struct {
	Name         string       `json:"name"`
	Abbreviation string       `json:"abbreviation"`
	Children     []GroupChild `json:"children"`
}

type GroupChild struct {
	Name  string       `json:"name"`
	Teams []GroupTeams `json:"teams"`
}

type GroupTeams struct {
	Id               string     `json:"id"`
	Slug             string     `json:"slug"`
	Name             string     `json:"name"`
	Abbreviation     string     `json:"abbreviation"`
	DisplayName      string     `json:"displayName"`
	ShortDisplayName string     `json:"shortDisplayName"`
	Logos            []TeamLogo `json:"logos"`
}

type TeamLogo struct {
	Href        string `json:"href"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Alt         string `json:"alt"`
	LastUpdated string `json:"lastUpdated"`
}

/* Season */

type ESPNSeason struct {
	Leagues []League      `json:"leagues"`
	Events  []SeasonEvent `json:"events"`
}

type League struct {
	Id           string           `json:"id"`
	Uid          string           `json:"uid"`
	Name         string           `json:"name"`
	Abbreviation string           `json:"abbreviation"`
	Slug         string           `json:"slug"`
	Season       SeasonInfo       `json:"season"`
	Calender     []SeasonCalendar `json:"calendar"`
}

type SeasonInfo struct {
	Year        int    `json:"year"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	DisplayName string `json:"DisplayName"`
}

type SeasonCalendar struct {
	Label     string          `json:"label"`
	Value     string          `json:"value"`
	StartDate string          `json:"startDate"`
	EndDate   string          `json:"endDate"`
	Entries   []CalenderEntry `json:"entries"`
}

type CalenderEntry struct {
	Label     string `json:"label"`
	AltLabel  string `json:"alternateLabel"`
	Detail    string `json:"detail"`
	Value     string `json:"value"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type SeasonEvent struct {
	Id           string              `json:"id"`
	Uid          string              `json:"uid"`
	Date         string              `json:"date"`
	Name         string              `json:"name"`
	ShortName    string              `json:"shortName"`
	Season       SeasonEventInfo     `json:"season"`
	Week         SeasonEventWeek     `json:"week"`
	Link         []SeasonEventLink   `json:"links"`
	Competitions []SeasonCompetition `json:"competitions"`
	// Status       SeasonEventStatus   `json:"status"`
}

type SeasonEventInfo struct {
	Year int `json:"year"`
}

type SeasonEventWeek struct {
	Number int `json:"number"`
}

type SeasonEventLink struct {
	Href      string `json:"href"`
	Text      string `json:"text"`
	ShortText string `json:"shortText"`
}

type SeasonCompetition struct {
	Id          string       `json:"id"`
	Uid         string       `json:"uid"`
	Date        string       `json:"date"`
	Competitors []Competitor `json:"competitors"`
}

type Competitor struct {
	Id       string         `json:"id"`
	Uid      string         `json:"uid"`
	HomeAway string         `json:"homeAway"`
	Winner   bool           `json:"winner"`
	Team     CompetitorTeam `json:"team"`
}

type CompetitorTeam struct {
	Id               string `json:"id"`
	Uid              string `json:"uid"`
	Location         string `json:"location"`
	Name             string `json:"name"`
	Abbreviation     string `json:"abbreviation"`
	DisplayName      string `json:"displayName"`
	ShortDisplayName string `json:"shortDisplayName"`
	Color            string `json:"color"`
	AltColor         string `json:"alternateColor"`
	ConferenceId     string `json:"conferenceId"`
}

type SeasonEventStatus struct {
	Clock        int             `json:"clock"`
	DisplayClock string          `json:"displayClock"`
	Period       int             `json:"period"`
	Type         EventStatusType `json:"type"`
}

type EventStatusType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
	ShortDetail string `json:"shortDetail"`
	AltDetail   string `json:"altDetail"`
}

/* Teams (All) */
type ESPNTeams struct {
	Sports []ESPNTeamsSport `json:"sports"`
}

type ESPNTeamsSport struct {
	Id      string                   `json:"id"`
	Leagues []ESPNTeamsSportsLeagues `json:"leagues"`
}

type ESPNTeamsSportsLeagues struct {
	Id    string                      `json:"id"`
	Teams []ESPNTeamsSportsLeagueTeam `json:"teams"`
}

type ESPNTeamsSportsLeagueTeam struct {
	Team struct {
		Id string `json:"id"`
	} `json:"team"`
}

/* Team (Individual against teamId) */
type ESPNTeamWrapper struct {
	Team ESPNCfbTeam `json:"team"`
}

type ESPNCfbTeam struct {
	Id               string      `json:"id"`
	Uid              string      `json:"uid"`
	Slug             string      `json:"slug"`
	Location         string      `json:"location"`
	Name             string      `json:"name"`
	Nickname         string      `json:"nickName"`
	Abbreviation     string      `json:"abbreviation"`
	DisplayName      string      `json:"displayName"`
	ShortDisplayName string      `json:"shortDisplayName"`
	Color            string      `json:"color"`
	AltColor         string      `json:"alternateColor"`
	IsActive         bool        `json:"isActive"`
	Logos            []ETeamLogo `json:"logos"`
	Groups           TeamGroups  `json:"groups"`
	Links            []ESPNLink  `json:"links"`
	StandingSummary  string      `json:"standingSummary"`
}

//TODO remove E prefix when old type are gone
type ETeamLogo struct {
	Href        string   `json:"href"`
	Width       int      `json:"width"`
	Heigt       int      `json:"height"`
	Alt         string   `json:"alt"`
	Rel         []string `json:"rel"`
	LastUpdated string   `json:"lastUpdated"`
}

type TeamGroups struct {
	Id     string `json:"id"`
	Parent struct {
		Id string `json:"id"`
	} `json:"parent"`
	IsConference bool `json:"isConference"`
}

/* Game */
type ESPNCfbGame struct {
	Boxscore GameBoxScore `json:"boxscore"`
	Header   GameHeader   `json:"header"`
	Links    []ESPNLink   `json:"links"`
	Week     int          `json:"week"`
}

type GameHeader struct {
	Id     string `json:"id"`
	Uid    string `json:"uid"`
	Season struct {
		Year int `json:"year"`
		Type int `json:"type"`
	} `json:"season"`
}

type GameBoxScore struct {
	Teams []BoxScoreTeam `json:"teams"`
}

type BoxScoreTeam struct {
	Id           string             `json:"id"`
	Statistics   []BoxScoreTeamStat `json:"statistics"`
	DisplayOrder int                `json:"displayOrder"`
	HomeAway     string             `json:"homeAway"`
}

type BoxScoreTeamStat struct {
	Name         string `json:"name"`
	DisplayValue string `json:"displayValue"`
	Label        string `json:"label"`
}

/* General */
type ESPNLink struct {
	Href      string   `json:"href"`
	Text      string   `json:"text"`
	ShortText string   `json:"shortText"`
	Rel       []string `json:"rel"`
}
