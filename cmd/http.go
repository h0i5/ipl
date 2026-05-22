package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetApiData is a blind caller + parser and uses generics
func GetApiData[T any](route string) (T, error) {
	const apiURL = "https://ipl-okn0.onrender.com"
	var result T
	resp, err := http.Get(apiURL + "/" + route)
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
