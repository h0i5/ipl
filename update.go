package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/h0i5/ipl/cmd"
)

func (m Model) handleTabNavigation(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "m":
		m.selectedTab = MatchView
	case "p":
		m.selectedTab = PointsTableView
	case "s":
		m.selectedTab = ScheduleView
	case "a":
		m.selectedTab = AboutView
	case "l":
		m.selectedTab = LiveView
	}
	m.currentView = m.selectedTab
	if m.currentView == MatchView {
		m.matchTable.Focus()
	} else {
		m.matchTable.Blur()
	}
	return m, nil
}

func (m Model) handleNavToTabView() (tea.Model, tea.Cmd) {
	m.selectedTab = m.currentView
	m.currentView = TabView
	m.matchTable.Blur()
	return m, nil
}

func (m Model) tabCursorIndex() int {
	for i, v := range tabOrder {
		if v == m.selectedTab {
			return i
		}
	}
	return 0
}

func (m Model) handleTabCursor(key string) (tea.Model, tea.Cmd) {
	idx := m.tabCursorIndex()
	switch key {
	case "up":
		idx--
	case "down":
		idx++
	case "enter", "right":
		m.currentView = m.selectedTab
		if m.currentView == MatchView {
			m.matchTable.Focus()
		}
		return m, nil
	}
	if idx < 0 {
		idx = len(tabOrder) - 1
	} else if idx >= len(tabOrder) {
		idx = 0
	}
	m.selectedTab = tabOrder[idx]
	return m, nil
}

func (m Model) handleQuit(key string) (tea.Model, tea.Cmd) {
	if key == "ctrl+c" || key == "q" {
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) allLoaded() bool {
	return !m.loadingMap[MatchView] &&
		!m.loadingMap[ScheduleView] &&
		!m.loadingMap[PointsTableView] &&
		!m.loadingMap[LiveView]
}

type squadMsg struct {
	slug string
	data cmd.SquadResponse
}

func fetchSquad(slug string) tea.Cmd {
	return func() tea.Msg {
		data, err := cmd.GetSquad(slug)
		if err != nil {
			return nil
		}
		return squadMsg{slug: slug, data: data}
	}
}

func squadCmdsForLive(live cmd.LiveMatchResponse) tea.Cmd {
	if live.LiveCount == 0 || len(live.LiveScore) == 0 {
		return nil
	}
	seen := map[string]bool{}
	var cmds []tea.Cmd
	for _, match := range live.LiveScore {
		for _, name := range []string{match.Team1, match.Team2} {
			slug := cmd.TeamToSlug(name)
			if slug != "" && !seen[slug] {
				seen[slug] = true
				cmds = append(cmds, fetchSquad(slug))
			}
		}
	}
	return tea.Batch(cmds...)
}

type refreshtickMsg time.Time

func refreshTickCmd() tea.Cmd {

	// refresh scores every 10 seconds
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return refreshtickMsg(t)
	})
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// set the data on initial update called by Init()
	case cmd.MatchScoresResponse:
		m.loadingMap[MatchView] = false
		m.items.matchScores = msg
		m.matchTable = m.buildMatchTable(m.matchTable, msg)
		if m.allLoaded() {
			m.currentView = LiveView
			m.selectedTab = LiveView
			m.matchTable.Focus()
		}

	case cmd.MatchScheduleResponse:
		m.loadingMap[ScheduleView] = false
		m.items.matchSchedule = msg
		if m.allLoaded() {
			m.currentView = LiveView
			m.selectedTab = LiveView
			m.matchTable.Focus()
		}

	case cmd.LiveMatchResponse:
		m.loadingMap[LiveView] = false
		m.items.liveMatch = msg
		if m.allLoaded() {
			m.currentView = LiveView
			m.selectedTab = LiveView
			m.matchTable.Focus()
		}
		return m, squadCmdsForLive(msg)

	case cmd.PointsTableResponse:
		m.loadingMap[PointsTableView] = false
		m.items.pointsTable = msg
		if m.allLoaded() {
			m.currentView = LiveView
			m.selectedTab = LiveView
			m.matchTable.Focus()
		}

	case tickMsg:
		m.showLoadingCursor = !m.showLoadingCursor
		if !m.allLoaded() {
			return m, tickCmd()
		}
		return m, nil

	case refreshtickMsg:
		var err error
		m.items.liveMatch, err = cmd.GetLiveMatchScores()
		if err != nil {
			return m, refreshTickCmd()
		}
		m.lastUpdated = time.Now()
		return m, tea.Batch(refreshTickCmd(), squadCmdsForLive(m.items.liveMatch))

	case squadMsg:
		m.items.squads[msg.slug] = msg.data

	case tea.KeyMsg:
		key := msg.String()

		if model, command := m.handleQuit(key); command != nil {
			return model, command
		}

		// letter shortcuts work from anywhere
		switch key {
		case "m", "p", "s", "a", "l":
			return m.handleTabNavigation(key)
		}

		// view-specific keys
		switch m.currentView {
		case MatchView:
			switch key {
			case "left":
				return m.handleNavToTabView()
			default:
				var command tea.Cmd
				m.matchTable, command = m.matchTable.Update(msg)
				return m, command
			}
		case PointsTableView, ScheduleView, AboutView, LiveView:
			if key == "left" {
				return m.handleNavToTabView()
			}
		case TabView:
			return m.handleTabCursor(key)
		case InitialLoadView:
			// ignore keys during initial load
		}

	// set the width and height to align content
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.matchTable.SetHeight(int(float64(m.height)*0.8) - 10)

	}

	return m, nil
}
