package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"pitwall/ui"
)

// rootCmd is the Cobra root command.
var rootCmd = &cobra.Command{
	Use:   "pitwall",
	Short: "Pitwall - Formula 1 race control TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		m := ui.NewModel()
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
