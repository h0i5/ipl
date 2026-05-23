package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/h0i5/ipl/cmd"
)

var tabOrder = []int{LiveView, MatchView, PointsTableView, ScheduleView, HistoricalView, AboutView}

var tabLabels = map[int]struct{ key, label string }{
	LiveView:        {"l", "live"},
	MatchView:       {"m", "matches"},
	PointsTableView: {"p", "points"},
	ScheduleView:    {"s", "schedule"},
	HistoricalView:  {"h", "historical"},
	AboutView:       {"a", "about"},
}

// clamp returns v clamped to [lo, hi]
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ── Root ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	var cursor string
	if m.showLoadingCursor {
		cursor = m.styles.cursorStyle.Render("█")
	} else {
		cursor = " "
	}

	if m.currentView == InitialLoadView {
		content := lipgloss.JoinHorizontal(
			lipgloss.Center,
			"IPL 2026",
			cursor,
		)

		return m.styles.loading.
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render(content)
	}

	// Cap the total UI size so it doesn't sprawl on huge terminals
	totalW := clamp(m.width-4, 80, 160)
	totalH := int(float64(m.height) * 0.9)

	sidebarW := 24
	// Body gets the rest; -2 accounts for the sidebar's right border
	bodyW := totalW - sidebarW - 2

	sidebar := m.renderSidebar(sidebarW, totalH)
	body := m.renderBody(bodyW, totalH)

	inner := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, body)

	// Single outer border around the whole thing
	framed := m.styles.outerBorder.
		Width(totalW).
		Render(inner)

	// Center in the actual terminal window
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		framed,
	)
}

// ── Sidebar ───────────────────────────────────────────────────────────────────

func (m Model) renderSidebar(width, height int) string {
	s := m.styles

	innerW := width - 1 // -1 for the right border

	titleText := lipgloss.JoinHorizontal(lipgloss.Center,
		s.gold.Bold(true).Render("IPL 2026"),
		s.muted.Render(" · h0i5"),
	)
	title := lipgloss.NewStyle().Width(innerW).Align(lipgloss.Center).Render(titleText)

	divider := s.faint.Render(strings.Repeat("─", innerW))

	inTabView := m.currentView == TabView

	var rows []string
	rows = append(rows, title, divider, "")

	for _, view := range tabOrder {
		info := tabLabels[view]
		keyStr := fmt.Sprintf("[%s]", info.key)

		isActive := m.currentView == view              // currently open view
		isCursor := inTabView && view == m.selectedTab // hovered in tab nav

		var row string
		switch {
		case isActive:
			key := s.tabKey.Render(keyStr)
			label := s.tabActive.Render(info.label)
			row = lipgloss.JoinHorizontal(lipgloss.Center, " ", key, label)

		case isCursor:
			key := s.tabKey.Render(keyStr)
			label := s.highlight.Padding(0, 2).Render(info.label)
			row = lipgloss.JoinHorizontal(lipgloss.Center, " ", key, label)

		default:
			key := s.tabKeyDim.Render(keyStr)
			label := s.tabInactive.Render(info.label)
			row = lipgloss.JoinHorizontal(lipgloss.Center, " ", key, label)
		}

		rows = append(rows, row)
	}

	rows = append(rows, "", divider, "")
	rows = append(rows, s.faint.Render("  ↑↓  navigate"))
	if m.currentView != TabView {
		rows = append(rows, s.faint.Render("  ←   focus"))
	} else {
		rows = append(rows, s.faint.Render("  →   move"))
	}
	rows = append(rows, s.faint.Render("  q   quit"))

	content := strings.Join(rows, "\n")

	return s.sidebar.
		Width(width).
		Height(height).
		Render(content)
}

// ── Body ──────────────────────────────────────────────────────────────────────

func (m Model) renderBody(width, height int) string {
	s := m.styles

	innerW := width - 4 // subtract body padding (2 each side)

	activeView := m.currentView
	if m.currentView == TabView {
		activeView = m.selectedTab
	}

	var content string
	switch activeView {
	case LiveView:
		content = m.renderLive(innerW)

	case MatchView:
		content = m.renderMatches(innerW)

	case PointsTableView:
		content = m.renderStandings(innerW)

	case ScheduleView:
		content = m.renderSchedule(innerW)

	case HistoricalView:
		content = m.renderHistorical(innerW)

	case AboutView:
		content = m.renderAbout(innerW)
	}

	return s.body.
		Width(width).
		Height(height).
		Render(content)
}

// ── Matches ───────────────────────────────────────────────────────────────────

func (m Model) renderMatches(width int) string {
	s := m.styles

	if m.loadingMap[MatchView] {
		return s.loading.Render("fetching matches...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	tbl := m.matchTable
	if m.currentView == TabView {
		ts := m.matchTableStyles
		ts.Selected = lipgloss.NewStyle()
		tbl.SetStyles(ts)
	}

	var sb strings.Builder
	sb.WriteString(heading.Render("all matches"))
	sb.WriteString("\n\n")
	sb.WriteString(tbl.View())
	sb.WriteString("\n")
	sb.WriteString(s.faint.Render("  ↑↓ scroll  ← back"))

	return sb.String()
}

func (m Model) renderMatchCard(match cmd.Match, width int, live bool) string {
	s := m.styles

	// Status badge
	var badge string
	if live {
		badge = s.liveDot.Render("● LIVE")
	} else {
		badge = s.faint.Render("✓ " + match.Status)
	}

	// Teams + scores row
	team1 := s.teamName.Render(match.Team1)
	team2 := s.teamName.Render(match.Team2)
	score1 := s.score.Render(match.Score1)
	score2 := s.score.Render(match.Score2)

	teamsRow := fmt.Sprintf("%s  %s    %s  %s",
		team1, score1, team2, score2,
	)

	// Venue + result
	venueStr := s.venue.Render("📍 " + match.Venue)

	var resultStr string
	if match.Result != "" {
		resultStr = s.result.Render("→ " + match.Result)
	}

	lines := []string{badge, teamsRow, venueStr}
	if resultStr != "" {
		lines = append(lines, resultStr)
	}
	cardW := width - 6
	if cardW > 60 {
		cardW = 60
	}
	return s.matchCard.Width(cardW).Render(strings.Join(lines, "\n"))
}

// ── Live ──────────────────────────────────────────────────────────────────────

func (m Model) renderLive(width int) string {
	s := m.styles

	if m.loadingMap[LiveView] {
		return s.loading.Render("Fetching live scores...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	var sb strings.Builder
	sb.WriteString(heading.Render("live matches"))
	sb.WriteString("\n\n")

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.UTC
	}

	updatedLine := s.faint.Render("last updated " + m.lastUpdated.In(loc).Format("15:04:05") + " • auto updates every 5s")

	data := m.items.liveMatch
	if data.LiveCount == 0 || len(data.LiveScore) == 0 {
		sb.WriteString(s.muted.Render("No matches live right now"))
		sb.WriteString(s.faint.Render("  •  " + m.lastUpdated.In(loc).Format("15:04:05")))
		sb.WriteString("\n")
		sb.WriteString(s.faint.Render("Press [m] for historical view."))
		return sb.String()
	}

	keys := make([]string, 0, len(data.LiveScore))
	for k := range data.LiveScore {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sb.WriteString(updatedLine)
	sb.WriteString("\n\n")

	for _, k := range keys {
		match := data.LiveScore[k]

		dot := "●"
		if !m.showLiveCursor {
			dot = " "
		}
		badge := s.liveDot.Render(dot + " LIVE")
		if d := match.LiveDetails; d != nil && d.MatchNumber != "" {
			badge = lipgloss.JoinHorizontal(lipgloss.Left,
				badge,
				s.faint.Render(fmt.Sprintf("  •  Match %s  •  Inning %d", d.MatchNumber, d.Inning)),
			)
		}
		sb.WriteString(badge + "\n\n")

		if width >= 90 {
			colW := width / 3

			// determine which innings is live
			liveScore, liveOvers := match.Score1, match.Overs1
			var targetRuns int
			hasTarget := false
			if d := match.LiveDetails; d != nil && d.Inning == 2 {
				liveScore, liveOvers = match.Score2, match.Overs2
				if parts := strings.SplitN(match.Score1, "/", 2); len(parts) > 0 {
					if r, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
						targetRuns = r + 1
						hasTarget = true
					}
				}
			} else if match.Score1 == "" || match.Score1 == "Yet to bat" {
				liveScore, liveOvers = match.Score2, match.Overs2
			}

			// build center block
			var midLines []string
			if art := scoreArt(liveScore); art != "" {
				midLines = append(midLines, art)
			} else {
				midLines = append(midLines, s.score.Render(liveScore))
			}
			if liveOvers != "" && liveOvers != "N.A" {
				midLines = append(midLines, s.faint.Render("("+liveOvers+" ov)"))
			}
			if hasTarget {
				midLines = append(midLines, "")
				midLines = append(midLines, s.muted.Render(fmt.Sprintf("Target: %d", targetRuns)))
			}
			if d := match.LiveDetails; d != nil {
				var rates []string
				if d.Rates.CRR != "" {
					rates = append(rates, "CRR "+d.Rates.CRR)
				}
				if d.Rates.RRR != "" && d.Rates.RRR != "--" {
					rates = append(rates, "RRR "+d.Rates.RRR)
				}
				if len(rates) > 0 {
					midLines = append(midLines, s.faint.Render(strings.Join(rates, "  •  ")))
				}
			}

			midW := width - 2*colW
			leftCol := lipgloss.NewStyle().Width(colW).Render(teamArt(match.Team1))
			midCol := lipgloss.NewStyle().Width(midW).Align(lipgloss.Center).Render(strings.Join(midLines, "\n"))
			rightCol := lipgloss.NewStyle().Width(colW).Align(lipgloss.Right).Render(teamArt(match.Team2))

			sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftCol, midCol, rightCol))
			sb.WriteString("\n")

			if d := match.LiveDetails; d != nil {
				if d.Venue != "" {
					sb.WriteString(s.venue.Render("📍 "+d.Venue) + "\n")
				}
				if d.Toss != "" {
					sb.WriteString(s.faint.Render("🪙 "+d.Toss) + "\n")
				}
			} else if match.StartTime != "" {
				displayTime := match.StartTime
				if t, err := time.Parse(time.RFC3339, match.StartTime); err == nil {
					displayTime = t.In(loc).Format("15:04 IST")
				}
				sb.WriteString(s.faint.Render("start  •  "+displayTime) + "\n")
			}
		} else {
			// narrow fallback
			overs1, overs2 := "", ""
			if match.Overs1 != "" && match.Overs1 != "N.A" {
				overs1 = "  " + s.faint.Render("("+match.Overs1+")")
			}
			if match.Overs2 != "" && match.Overs2 != "N.A" {
				overs2 = "  " + s.faint.Render("("+match.Overs2+")")
			}
			row1 := lipgloss.JoinHorizontal(lipgloss.Left,
				s.teamName.Render(match.Team1), "  ", s.score.Render(match.Score1), overs1,
			)
			row2 := lipgloss.JoinHorizontal(lipgloss.Left,
				s.teamName.Render(match.Team2), "  ", s.score.Render(match.Score2), overs2,
			)
			var lines []string
			lines = append(lines, "", row1, row2)
			if d := match.LiveDetails; d != nil {
				if d.Venue != "" {
					lines = append(lines, "", s.venue.Render("📍 "+d.Venue))
				}
				if d.Toss != "" {
					lines = append(lines, s.faint.Render("🪙 "+d.Toss))
				}
			} else if match.StartTime != "" {
				displayTime := match.StartTime
				if t, err := time.Parse(time.RFC3339, match.StartTime); err == nil {
					displayTime = t.In(loc).Format("15:04 IST")
				}
				lines = append(lines, "", s.faint.Render("start  •  "+displayTime))
			}
			narrowW := width - 6
			if narrowW > 60 {
				narrowW = 60
			}
			sb.WriteString(s.matchCard.Width(narrowW).Render(strings.Join(lines, "\n")))
		}

		sb.WriteString("\n\n")

		if match.LiveDetails != nil {
			sb.WriteString(m.renderLiveDetails(match.LiveDetails, width))
			sb.WriteString("\n\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m Model) renderLiveDetails(d *cmd.LiveDetails, width int) string {
	s := m.styles
	var sb strings.Builder

	sec := func(title string) string {
		return s.muted.Render(title)
	}

	// batters
	if len(d.Batters) > 0 {
		sb.WriteString(sec("BATTERS") + "\n")
		for _, b := range d.Batters {
			balls := strings.Trim(b.Balls, "()")
			prefix := "  "
			if b.OnStrike {
				prefix = s.gold.Render("* ")
			}
			line := fmt.Sprintf("%-22s  %3s (%s)  4s:%-2s 6s:%-2s  SR:%s",
				b.FullName, b.Runs, balls, b.Fours, b.Sixes, b.StrikeRate)
			sb.WriteString(prefix + s.score.Render(line) + "\n")
		}
		sb.WriteString("\n")
	}

	// bowler
	if d.Bowler.FullName != "" {
		sb.WriteString(sec("BOWLER") + "\n")
		line := fmt.Sprintf("  %-22s  %s  (%s ov)  Econ: %s",
			d.Bowler.FullName, d.Bowler.Figures, d.Bowler.Overs, d.Bowler.Economy)
		sb.WriteString(s.muted.Render(line) + "\n\n")
	}

	// partnership + rates on one line
	pLine := fmt.Sprintf("Partnership: %d runs  %d balls", d.Partnership.Runs, d.Partnership.Balls)
	if d.Rates.CRR != "" {
		pLine += fmt.Sprintf("    CRR: %s", d.Rates.CRR)
	}
	if d.Rates.RRR != "" && d.Rates.RRR != "--" {
		pLine += fmt.Sprintf("  RRR: %s", d.Rates.RRR)
	}
	sb.WriteString(s.faint.Render(pLine) + "\n")

	// last wicket
	if d.LastWicket.Name != "" {
		balls := strings.Trim(d.LastWicket.Balls, "()")
		sb.WriteString(s.faint.Render(fmt.Sprintf("Last wicket: %s  %s(%s)", d.LastWicket.Name, d.LastWicket.Runs, balls)) + "\n")
	}

	// recent overs (latest last → show last 4, reversed so latest is on top)
	if len(d.RecentOvers) > 0 {
		sb.WriteString("\n" + sec("RECENT OVERS") + "\n")
		overs := d.RecentOvers
		if len(overs) > 4 {
			overs = overs[len(overs)-4:]
		}
		for i := len(overs) - 1; i >= 0; i-- {
			o := overs[i]
			var balls strings.Builder
			for j, ball := range o.OverInfo {
				if j > 0 {
					balls.WriteString("  ")
				}
				balls.WriteString(renderBall(s, ball))
			}
			sb.WriteString(fmt.Sprintf("  %-8s  %s   %s\n",
				o.Over, balls.String(), s.faint.Render(fmt.Sprintf("[%d]", o.Total))))
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}

func renderBall(s Styles, ball string) string {
	switch ball {
	case "W":
		return s.liveDot.Render(ball)
	case "4":
		return s.gold.Render(ball)
	case "6":
		return s.result.Render(ball)
	case "wd", "nb", "lb", "b":
		return s.muted.Render(ball)
	case "0":
		return s.faint.Render(ball)
	default:
		return s.score.Render(ball)
	}
}

func (m Model) renderSquads(slug1, slug2, name1, name2 string, width int) string {
	s := m.styles
	sq1, ok1 := m.items.squads[cmd.TeamToSlug(slug1)]
	sq2, ok2 := m.items.squads[cmd.TeamToSlug(slug2)]

	if !ok1 && !ok2 {
		return s.faint.Render("  loading squads...")
	}

	colW := (width - 4) / 2
	if colW < 20 {
		colW = 20
	}

	left := buildSquadColumn(s, name1, sq1, ok1, colW)
	right := buildSquadColumn(s, name2, sq2, ok2, colW)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(colW).Render(left),
		"    ",
		lipgloss.NewStyle().Width(colW).Render(right),
	)
}

func buildSquadColumn(s Styles, teamName string, sq cmd.SquadResponse, loaded bool, width int) string {
	divider := s.faint.Render(strings.Repeat("─", width-2))
	header := s.gold.Bold(true).Render(teamName)

	if !loaded {
		return strings.Join([]string{header, divider, s.faint.Render("  loading...")}, "\n")
	}

	type entry struct {
		key string
		p   cmd.SquadPlayer
	}
	seen := map[string]bool{}
	players := make([]entry, 0, len(sq.Squad))
	for k, p := range sq.Squad {
		if seen[p.Name] {
			continue
		}
		seen[p.Name] = true
		players = append(players, entry{k, p})
	}
	sort.Slice(players, func(i, j int) bool {
		ri, rj := squadPlayerRank(players[i].p), squadPlayerRank(players[j].p)
		if ri != rj {
			return ri < rj
		}
		return players[i].p.Name < players[j].p.Name
	})

	lines := []string{header, divider}
	for _, e := range players {
		p := e.p
		name := p.Name
		if len([]rune(name)) > 22 {
			name = string([]rune(name)[:21]) + "…"
		}
		nameW := lipgloss.NewStyle().Width(24).Render(s.muted.Render(name))
		tag := squadPlayerTag(p)
		lines = append(lines, "  "+nameW+s.faint.Render(tag))
	}

	return strings.Join(lines, "\n")
}

func squadPlayerRank(p cmd.SquadPlayer) int {
	if p.Wicketkeeper {
		return 0
	}
	if strings.Contains(p.Style, "Bat") {
		return 1
	}
	return 2
}

func squadPlayerTag(p cmd.SquadPlayer) string {
	var tag string
	switch {
	case p.Wicketkeeper:
		tag = "WK  "
	case strings.Contains(p.Style, "Bat"):
		tag = "BAT "
	default:
		tag = "BOWL"
	}
	if p.Overseas {
		tag += " OS"
	}
	return tag
}

// ── Standings ─────────────────────────────────────────────────────────────────

func (m Model) renderStandings(width int) string {
	s := m.styles

	if m.loadingMap[PointsTableView] {
		return s.loading.Render("fetching standings...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	data := m.items.pointsTable
	var sb strings.Builder

	sb.WriteString(heading.Render("points table"))
	sb.WriteString("\n\n")

	// Header row
	header := fmt.Sprintf("%-3s %-29s %3s %3s %3s %3s %7s %4s",
		"#", "Team", "P", "W", "L", "NR", "NRR", "Pts",
	)
	sb.WriteString(s.tableHeader.Render(header))
	sb.WriteString("\n")
	sb.WriteString(s.faint.Render(strings.Repeat("─", width)))
	sb.WriteString("\n")

	keys := make([]string, 0, len(data.PointsTable))
	for k := range data.PointsTable {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := data.PointsTable[keys[i]], data.PointsTable[keys[j]]
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		return a.NetRunRate > b.NetRunRate
	})

	for i, k := range keys {
		team := data.PointsTable[k]
		row := fmt.Sprintf("%-3d %-29s %3d %3d %3d %3d %+7.3f %4d",
			i+1,
			team.Name,
			team.Played,
			team.Won,
			team.Loss,
			team.NoResult,
			team.NetRunRate,
			team.Points,
		)

		// Top 4 qualify for playoffs — render brighter
		if i < 4 {
			sb.WriteString(s.tableRow.Render(row))
		} else {
			sb.WriteString(s.tableRowAlt.Render(row))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(s.faint.Render("top 4 qualify for playoffs"))

	return sb.String()
}

// ── Schedule ──────────────────────────────────────────────────────────────────

func (m Model) renderSchedule(width int) string {
	s := m.styles

	if m.loadingMap[ScheduleView] {
		return s.loading.Render("fetching schedule...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	data := m.items.matchSchedule
	var sb strings.Builder

	sb.WriteString(heading.Render("upcoming matches"))
	sb.WriteString("\n\n")

	keys := make([]string, 0, len(data.Schedule))
	for k := range data.Schedule {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	header := fmt.Sprintf("%-8s  %-32s %-18s %-8s %s", "#", "Teams", "Date", "Time", "Venue")
	sb.WriteString(s.tableHeader.Render(header))
	sb.WriteString("\n")
	sb.WriteString(s.faint.Render(strings.Repeat("─", width)))
	sb.WriteString("\n")

	for i, k := range keys {
		match := data.Schedule[k]
		row := fmt.Sprintf("%-8s  %-32s %-18s %-8s %s",
			k,
			truncate(match.Rival, 31),
			truncate(match.Date, 17),
			match.Time,
			truncateVenue(match.Location, 35),
		)
		if i%2 == 0 {
			sb.WriteString(s.tableRow.Render(row))
		} else {
			sb.WriteString(s.tableRowAlt.Render(row))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ── Historical ────────────────────────────────────────────────────────────────

func (m Model) renderHistorical(width int) string {
	s := m.styles

	if m.loadingMap[HistoricalView] {
		return s.loading.Render("fetching winners...")
	}

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	data := m.items.historicalWinners
	years := make([]string, 0, len(data.Winners))
	for y := range data.Winners {
		years = append(years, y)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(years)))

	var sb strings.Builder
	sb.WriteString(heading.Render("ipl winners"))
	sb.WriteString("\n\n")

	header := fmt.Sprintf("%-6s  %-34s %-20s %-14s %s", "Year", "Winner", "Runner Up", "Venue", "Result")
	sb.WriteString(s.tableHeader.Render(header))
	sb.WriteString("\n")
	sb.WriteString(s.faint.Render(strings.Repeat("─", width)))
	sb.WriteString("\n")

	for i, yr := range years {
		w := data.Winners[yr]
		row := fmt.Sprintf("%-6s  %-34s %-20s %-14s %s",
			yr,
			truncate(w.Winner, 33),
			truncate(w.RunnerUp, 19),
			truncate(w.Venue, 13),
			w.WonBy,
		)
		if i%2 == 0 {
			sb.WriteString(s.tableRow.Render(row))
		} else {
			sb.WriteString(s.tableRowAlt.Render(row))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	if m.currentView == HistoricalView {
		sb.WriteString(s.faint.Render("← back"))
	}

	return sb.String()
}

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}

func truncateVenue(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	words := strings.Fields(s)
	if len(words) <= 1 {
		return truncate(s, max)
	}
	last := words[len(words)-1]
	prefix := words[:len(words)-1]
	for len(prefix) > 0 {
		candidate := strings.Join(prefix, " ") + "… " + last
		if len([]rune(candidate)) <= max {
			return candidate
		}
		prefix = prefix[:len(prefix)-1]
	}
	return "… " + last
}

// ── About ─────────────────────────────────────────────────────────────────────

func (m Model) renderAbout(width int) string {
	s := m.styles

	heading := s.header
	if m.currentView == TabView {
		heading = s.muted
	}

	lines := []string{
		heading.Bold(true).Render("ipl tui"),
		"",
		s.muted.Render("a terminal viewer for ipl 2026."),
		"",
		s.faint.Render("why?"),
		s.faint.Render("idk"),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left,
			s.faint.Render("check out "),
			s.highlight.Render("harshiyer.in"),
			s.faint.Render(" :)"),
		),

		"",
		lipgloss.JoinHorizontal(lipgloss.Left,
			s.faint.Render("built with bubbletea • code @ "),
			s.highlight.Render("https://github.com/h0i5/ipl"),
		),
		"",
	}

	return strings.Join(lines, "\n")
}
