package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// GetApiData is a blind caller + parser and uses generics
func GetApiData[T any](route string) (T, error) {
	var result T
	apiBase := os.Getenv("API_URL")
	resp, err := http.Get(apiBase + "/" + route)
	if err != nil {
		return result, fmt.Errorf("error fetching %s: %w", route, err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("error decoding %s: %w", route, err)
	}
	return result, nil
}

func GetMatchScores() (MatchScoresResponse, error) {
	return GetApiData[MatchScoresResponse]("ipl-2026-live-score")
}

func GetMatchSchedule() (MatchScheduleResponse, error) {
	return GetApiData[MatchScheduleResponse]("ipl-2026-schedule")
}

func GetPointsTable() (PointsTableResponse, error) {
	return GetApiData[PointsTableResponse]("ipl-2026-points-table")
}

func GetLiveMatchScores() (LiveMatchResponse, error) {
	return GetApiData[LiveMatchResponse]("/ipl-2026-live-score-s2")
}

func GetHistoricalWinners() (HistoricalWinnersResponse, error) {
	return GetApiData[HistoricalWinnersResponse]("ipl-winners")
}

var teamSlugs = map[string]string{
	"mi": "mi", "mumbai indians": "mi", "mumbai": "mi",
	"csk": "csk", "chennai super kings": "csk", "chennai": "csk",
	"rcb": "rcb", "royal challengers bengaluru": "rcb", "royal challengers bangalore": "rcb", "royal challengers": "rcb",
	"kkr": "kkr", "kolkata knight riders": "kkr", "kolkata": "kkr",
	"srh": "srh", "sunrisers hyderabad": "srh", "sunrisers": "srh",
	"dc": "dc", "delhi capitals": "dc", "delhi": "dc",
	"pbks": "pbks", "punjab kings": "pbks", "punjab": "pbks",
	"rr": "rr", "rajasthan royals": "rr", "rajasthan": "rr",
	"gt": "gt", "gujarat titans": "gt", "gujarat": "gt",
	"lsg": "lsg", "lucknow super giants": "lsg", "lucknow": "lsg",
}

// keyword fallbacks for partial team names from the live API
var teamKeywords = []struct {
	keyword string
	slug    string
}{
	{"mumbai", "mi"},
	{"chennai", "csk"},
	{"challengers", "rcb"},
	{"kolkata", "kkr"},
	{"sunrisers", "srh"},
	{"delhi", "dc"},
	{"punjab", "pbks"},
	{"rajasthan", "rr"},
	{"gujarat", "gt"},
	{"lucknow", "lsg"},
}

func TeamToSlug(name string) string {
	lower := strings.ToLower(strings.TrimSpace(name))
	if slug, ok := teamSlugs[lower]; ok {
		return slug
	}
	for _, kw := range teamKeywords {
		if strings.Contains(lower, kw.keyword) {
			return kw.slug
		}
	}
	return ""
}

func GetSquad(slug string) (SquadResponse, error) {
	return GetApiData[SquadResponse]("squad/" + slug)
}
