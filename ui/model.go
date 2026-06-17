package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"pitwall/internal/models"
)

// FetchFunc fetches live domain data. Implemented by the application wiring (cmd).
type FetchFunc func(ctx context.Context) ([]models.Driver, []models.Event, models.Race, error)

// Model is the Bubble Tea model for the TUI.
type Model struct {
	Leaderboard []models.Driver
	Events      []models.Event
	NextRace    models.Race

	// live refresh
	fetchFn        FetchFunc
	refreshInterval time.Duration
	Loading        bool
	Err            string
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
		Loading:     false,
		Err:         "",
	}
}

// NewLiveModel constructs a Model that will auto-refresh using fetchFn every interval.
func NewLiveModel(fetchFn FetchFunc, interval time.Duration, leaderboard []models.Driver, events []models.Event, next models.Race) *Model {
	m := NewModelWithData(leaderboard, events, next)
	m.fetchFn = fetchFn
	m.refreshInterval = interval
	return m
}

// Messages
type tickMsg time.Time
type dataMsg struct {
	leaders []models.Driver
	events  []models.Event
	next    models.Race
}

type errMsg struct{ err error }

// Bubble Tea methods
func (m *Model) Init() tea.Cmd {
	// start with immediate fetch if fetchFn provided
	if m.fetchFn != nil {
		return tea.Batch(m.fetchOnceCmd(), tea.Tick(m.refreshInterval, func(t time.Time) tea.Msg { return tickMsg(t) }))
	}
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		if m.fetchFn != nil {
			return m, tea.Batch(m.fetchOnceCmd(), tea.Tick(m.refreshInterval, func(t time.Time) tea.Msg { return tickMsg(t) }))
		}
	case dataMsg:
		m.Leaderboard = msg.leaders
		m.Events = msg.events
		m.NextRace = msg.next
		m.Loading = false
		m.Err = ""
	case errMsg:
		m.Loading = false
		m.Err = msg.err.Error()
	}
	return m, nil
}

func (m *Model) View() string {
	return Render(m)
}

// fetchOnceCmd performs a single fetch using the configured fetchFn and returns a dataMsg or errMsg.
func (m *Model) fetchOnceCmd() tea.Cmd {
	return func() tea.Msg {
		m.Loading = true
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		leaders, events, next, err := m.fetchFn(ctx)
		if err != nil {
			return errMsg{err}
		}
		return dataMsg{leaders: leaders, events: events, next: next}
	}
}
