package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Render composes the entire UI from sections.
func Render(m *Model) string {
	header := titleStyle.Render(" Pitwall ")

	// Leaderboard
	var lbLines []string
	for _, d := range m.Leaderboard {
		lbLines = append(lbLines, fmt.Sprintf("%2d %-18s %6s", d.Position, d.Name, d.Gap))
	}
	leaderboard := sectionStyle.Copy().Width(40).Render(lipgloss.JoinVertical(lipgloss.Left, append([]string{"Leaderboard", ""}, lbLines...)...))

	// Events
	var evLines []string
	for _, e := range m.Events {
		evLines = append(evLines, eventStyle.Render(e.Message))
	}
	events := sectionStyle.Copy().Width(30).Render(lipgloss.JoinVertical(lipgloss.Left, append([]string{"Recent Events", ""}, evLines...)...))

	// Next Race
	next := sectionStyle.Copy().Width(30).Render(lipgloss.JoinVertical(lipgloss.Left, "Next Race", "", raceStyle.Render(m.NextRace.Name), m.NextRace.Date))

	// Layout: header + three columns
	columns := lipgloss.JoinHorizontal(lipgloss.Top, leaderboard, events, next)

	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("Press q to quit")

	content := lipgloss.JoinVertical(lipgloss.Center, header, "", columns, "", footer)

	return lipgloss.NewStyle().Align(lipgloss.Center).Render(content)
}
