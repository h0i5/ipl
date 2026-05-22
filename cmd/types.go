package cmd

type MatchScoresResponse struct {
	StatusCode  int              `json:"status_code"`
	Season      string           `json:"season"`
	Source      string           `json:"source"`
	Status      string           `json:"status"`
	LiveCount   int              `json:"live_count"`
	LiveScore   map[string]Match `json:"live_score"`
	Matches     map[string]Match `json:"matches"`
	DateChecked string           `json:"date_checked"`
}

type Match struct {
	Status       string `json:"status"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Venue        string `json:"venue"`
	Date         string `json:"date"`
	StartTimeUTC string `json:"start_time_utc"`
	Team1        string `json:"team_1"`
	Score1       string `json:"score_1"`
	Team2        string `json:"team_2"`
	Score2       string `json:"score_2"`
	Result       string `json:"result,omitempty"`
}

type MatchScheduleResponse struct {
	StatusCode int                       `json:"status_code"`
	Season     string                    `json:"season"`
	Schedule   map[string]ScheduledMatch `json:"schedule"`
}

type ScheduledMatch struct {
	Rival    string `json:"Rival"`
	Location string `json:"Location"`
	Date     string `json:"Date"`
	Time     string `json:"Time"`
}

type PointsTableResponse struct {
	StatusCode  int                     `json:"status_code"`
	Season      string                  `json:"season"`
	PointsTable map[string]TeamStanding `json:"points_table"`
}

type TeamStanding struct {
	Name       string  `json:"Name"`
	Played     int     `json:"Played"`
	Won        int     `json:"Won"`
	Loss       int     `json:"Loss"`
	NoResult   int     `json:"No Result"`
	NetRunRate float64 `json:"Net Run Rate"`
	Points     int     `json:"Points"`
}

type HistoricalWinnersResponse struct {
	StatusCode int                   `json:"status_code"`
	Winners    map[string]YearWinner `json:"winners"`
}

type YearWinner struct {
	Winner   string `json:"Winner"`
	WonBy    string `json:"Won By"`
	RunnerUp string `json:"Runner Up"`
	Venue    string `json:"Venue"`
}

type LiveMatch struct {
	Status    string `json:"status"`
	Team1     string `json:"team_1"`
	Score1    string `json:"score_1"`
	Overs1    string `json:"overs_1"`
	Team2     string `json:"team_2"`
	Score2    string `json:"score_2"`
	Overs2    string `json:"overs_2"`
	MatchURL  string `json:"match_url"`
	StartTime string `json:"start_time,omitempty"`
}

type LiveMatchResponse struct {
	StatusCode  int                  `json:"status_code"`
	Season      string               `json:"season"`
	Source      string               `json:"source"`
	Status      string               `json:"status"`
	LiveCount   int                  `json:"live_count"`
	LiveScore   map[string]LiveMatch `json:"live_score"`
	Matches     map[string]LiveMatch `json:"matches"`
	DateChecked string               `json:"date_checked"`
}

type SquadPlayer struct {
	Name         string `json:"Name"`
	Nationality  string `json:"Nationality"`
	Role         string `json:"Role"`
	Style        string `json:"Style"`
	Wicketkeeper bool   `json:"Wicketkeeper"`
	Overseas     bool   `json:"Overseas"`
}

type SquadResponse struct {
	StatusCode int                    `json:"status_code"`
	Team       string                 `json:"team"`
	Season     string                 `json:"season"`
	Squad      map[string]SquadPlayer `json:"squad"`
}

//type LiveMatchResponse struct {
//	StatusCode  int                    `json:"status_code"`
//	Season      string                 `json:"season"`
//	Source      string                 `json:"source"`
//	SeriesID    string                 `json:"series_id,omitempty"`
//	Status      string                 `json:"status"`
//	LiveCount   int                    `json:"live_count"`
//	LiveScore   map[string]interface{} `json:"live_score,omitempty"`
//	Matches     map[string]LiveMatch   `json:"matches"`
//	DateChecked string                 `json:"date_checked,omitempty"`
//}
//
//type LiveMatch struct {
//	Status     string `json:"status"`
//	Title      string `json:"title"`
//	Info       string `json:"info,omitempty"`
//	Team1      string `json:"team_1"`
//	Score1     string `json:"score_1"`
//	Team2      string `json:"team_2"`
//	Score2     string `json:"score_2"`
//	MatchURL   string `json:"match_url,omitempty"`
//	StatusText string `json:"status_text"`
//}
