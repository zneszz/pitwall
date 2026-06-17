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

	// Layout: header + three columns
	columns := lipgloss.JoinHorizontal(lipgloss.Top, leaderboard, events, next)

	footer := footerStyle.Render("q: quit  •  auto-refresh: 15s")

	content := lipgloss.JoinVertical(lipgloss.Center, header, "", columns, "", footer)

	return lipgloss.NewStyle().Align(lipgloss.Center).Render(content)
}
