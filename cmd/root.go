package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// rootCmd is the Cobra root command.
var rootCmd = &cobra.Command{
	Use:   "pitwall",
	Short: "Pitwall - Formula 1 race control TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create a simple Bubble Tea model that renders static text for now.
		m := &model{}
		p := tea.NewProgram(m)
		if err := p.Start(); err != nil {
			return fmt.Errorf("failed to start TUI: %w", err)
		}
		return nil
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// Minimal Bubble Tea model for the initial screen.
type model struct{}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	header := titleStyle.Render("Pitwall")
	sub := lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Render("Formula 1 Race Control")

	return lipgloss.NewStyle().Align(lipgloss.Center).Render(fmt.Sprintf("%s\n\n%s", header, sub))
}
