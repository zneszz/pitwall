package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Background(lipgloss.Color("160")).Padding(0, 1)
	sectionStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).Margin(0, 1)
	leaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	eventStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	raceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Italic(true)
)
