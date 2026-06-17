package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Header/title
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("160")).Padding(0, 2)
	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Padding(0, 1)

	// Section base style (dark background, subtle border)
	sectionStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2).Margin(0, 1).Background(lipgloss.Color("235"))

	// Column specific tweaks
	leaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	raceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("121")).Bold(true)

	// Footer
	footerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// helper functions for event line styles
func eventLine(msg string) string {
	if msg == "" {
		return ""
	}
	// pick color by keyword
	if contains(msg, "Yellow") || contains(msg, "🟡") {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render(msg)
	}
	if contains(msg, "Safety") || contains(msg, "🚗") {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Bold(true).Render(msg)
	}
	if contains(msg, "Green") || contains(msg, "🟢") {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render(msg)
	}
	// default
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render(msg)
}

// contains is a small wrapper around strings.Contains for readability.
func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
