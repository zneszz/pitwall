package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"pitwall/internal/models"
)

// Model is the Bubble Tea model for the TUI.
type Model struct {
	Leaderboard []models.Driver
	Events      []models.Event
	NextRace    models.Race
}

// NewModel returns a Model populated with mock data.
func NewModel() *Model {
	return NewModelWithData(
		[]models.Driver{
			{Position: 1, Name: "Max Verstappen", Gap: "+0.0"},
			{Position: 2, Name: "Lewis Hamilton", Gap: "+1.2"},
			{Position: 3, Name: "Charles Leclerc", Gap: "+2.5"},
		},
		[]models.Event{
			{Message: "🟡 Yellow Flag"},
			{Message: "🚗 Safety Car"},
			{Message: "🟢 Green Flag"},
		},
		models.Race{Name: "Monaco Grand Prix", Date: "2026-05-31"},
	)
}

// NewModelWithData constructs a UI model from provided domain data.
func NewModelWithData(leaderboard []models.Driver, events []models.Event, next models.Race) *Model {
	return &Model{
		Leaderboard: leaderboard,
		Events:      events,
		NextRace:    next,
	}
}

// Bubble Tea methods
func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() string {
	return Render(m)
}
