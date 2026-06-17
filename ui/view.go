package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Render composes the entire UI from sections.
func Render(m *Model) string {
	header := lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render(" Pitwall "), headerStyle.Render("Formula 1 Race Control"))

	// Leaderboard (compact)
	var lbLines []string
	for _, d := range m.Leaderboard {
		lbLines = append(lbLines, fmt.Sprintf("%2d %-14s %5s", d.Position, d.Name, d.Gap))
	}
	leaderboardBody := lipgloss.JoinVertical(lipgloss.Left, append([]string{"Leaderboard", ""}, lbLines...)...)
	leaderboard := sectionStyle.Copy().Width(36).Render(leaderboardBody)

	// Events
	var evLines []string
	for _, e := range m.Events {
		evLines = append(evLines, eventLine(e.Message))
	}
	eventsBody := lipgloss.JoinVertical(lipgloss.Left, append([]string{"Recent Events", ""}, evLines...)...)
	events := sectionStyle.Copy().Width(28).Render(eventsBody)

	// Next Race
	nextBody := lipgloss.JoinVertical(lipgloss.Left, "Next Race", "", raceStyle.Render(m.NextRace.Name), m.NextRace.Date)
	next := sectionStyle.Copy().Width(28).Render(nextBody)

	// Additional panels: Past leaders and Standings
	// Past leaders (compact)
	var leadersLines []string
	if len(m.PastLeaders) == 0 {
		leadersLines = []string{"No data"}
	} else {
		for i, l := range m.PastLeaders {
			leadersLines = append(leadersLines, fmt.Sprintf("%d. %s", i+1, l.Name))
		}
	}
	leadersPanel := sectionStyle.Copy().Width(36).Render(lipgloss.JoinVertical(lipgloss.Left, append([]string{"Past GP Leaders", ""}, leadersLines...)...))

	// Standings (drivers top 5 + teams top 3)
	var dsLines []string
	for i, ds := range m.DriverStandings {
		if i >= 5 { break }
		dsLines = append(dsLines, fmt.Sprintf("%2d %-14s %5.1f", ds.Position, ds.DriverName, ds.Points))
	}
	if len(dsLines) == 0 { dsLines = []string{"No data"} }
	driversPanel := sectionStyle.Copy().Width(36).Render(lipgloss.JoinVertical(lipgloss.Left, append([]string{"Driver Standings", ""}, dsLines...)...))

	var tsLines []string
	for i, ts := range m.TeamStandings {
		if i >= 3 { break }
		tsLines = append(tsLines, fmt.Sprintf("%2d %-14s %5.1f", ts.Position, ts.TeamName, ts.Points))
	}
	if len(tsLines) == 0 { tsLines = []string{"No data"} }
	teamsPanel := sectionStyle.Copy().Width(36).Render(lipgloss.JoinVertical(lipgloss.Left, append([]string{"Team Standings", ""}, tsLines...)...))

	// Layout: header + three columns + a second row with leaders and standings
	columns := lipgloss.JoinHorizontal(lipgloss.Top, leaderboard, events, next)
	second := lipgloss.JoinHorizontal(lipgloss.Top, leadersPanel, driversPanel, teamsPanel)

	footer := footerStyle.Render("q: quit  •  auto-refresh: 15s")

	content := lipgloss.JoinVertical(lipgloss.Center, header, "", columns, "", second, "", footer)

	return lipgloss.NewStyle().Align(lipgloss.Center).Render(content)
}
