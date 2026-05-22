package main

import (
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/h0i5/ipl/cmd"
)

func (model Model) buildMatchTable(t table.Model, data cmd.MatchScoresResponse) table.Model {
	type entry struct {
		match cmd.Match
		time  time.Time
		key   string
	}

	var live, upcoming, completed []entry

	for k, m := range data.Matches {
		ts, _ := time.Parse(time.RFC3339, m.StartTimeUTC)
		e := entry{match: m, time: ts, key: k}
		switch m.Status {
		case "Live":
			live = append(live, e)
		case "Upcoming":
			upcoming = append(upcoming, e)
		default: // "Completed"
			completed = append(completed, e)
		}
	}

	// Live: no sort needed, usually 1-2 matches
	// Upcoming: soonest first
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].time.Before(upcoming[j].time)
	})
	// Completed: most recent first
	sort.Slice(completed, func(i, j int) bool {
		return completed[i].time.After(completed[j].time)
	})

	// Combine: live → upcoming → completed
	ordered := append(live, append(upcoming, completed...)...)

	var rows []table.Row
	for _, e := range ordered {
		m := e.match

		var status string
		switch m.Status {
		case "Live":
			if model.showLiveCursor {
				status = "● Live"
			} else {
				status = "  Live"
			}
		case "Upcoming":
			status = "○ Soon"
		default:
			status = "✓ Done"
		}

		score1 := m.Score1
		score2 := m.Score2
		if score1 == "N.A" {
			score1 = "—"
		}
		if score2 == "N.A" {
			score2 = "—"
		}

		result := m.Result
		if result == "" {
			result = "—"
		}

		date := m.Date
		if date == "TBD" {
			date = e.time.Format("2 Jan")
		}

		rows = append(rows, table.Row{
			status,
			m.Team1 + " v " + m.Team2,
			score1,
			score2,
			m.Venue,
			result,
			date,
		})
	}

	t.SetRows(rows)
	return t
}
