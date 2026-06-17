package cmd

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"pitwall/internal/models"
	"pitwall/internal/api"
	"pitwall/services"
	"pitwall/ui"
)

// rootCmd is the Cobra root command.
var rootCmd = &cobra.Command{
	Use:   "pitwall",
	Short: "Pitwall - Formula 1 race control TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create API client and services
		client, err := api.NewClient("", nil)
		if err != nil {
			return fmt.Errorf("creating api client: %w", err)
		}
		eService := services.NewEventsService(client)
		sService := services.NewScheduleService(client)

		// Fetch data with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		events, err := eService.GetRecentEvents(ctx)
		if err != nil || len(events) == 0 {
			// Fallback to example events
			events = []models.Event{
				{Message: "🟡 Yellow Flag"},
				{Message: "🚗 Safety Car"},
				{Message: "🟢 Green Flag"},
			}
		}

		race, err := sService.GetNextRace(ctx)
		if err != nil {
			// fallback to placeholder
			race = models.Race{Name: "TBD", Date: ""}
		}

		// Leaderboard remains mocked for now
		leaderboard := []models.Driver{
			{Position: 1, Name: "Max Verstappen", Gap: "+0.0"},
			{Position: 2, Name: "Lewis Hamilton", Gap: "+1.2"},
			{Position: 3, Name: "Charles Leclerc", Gap: "+2.5"},
		}

		// build a fetch function that the UI can use for auto-refresh
		fetchFn := func(ctx context.Context) ([]models.Driver, []models.Event, models.Race, error) {
			eventsLive, err := eService.GetRecentEvents(ctx)
			if err != nil || len(eventsLive) == 0 {
				// fallback to static events when API fails
				eventsLive = []models.Event{
					{Message: "🟡 Yellow Flag"},
					{Message: "🚗 Safety Car"},
					{Message: "🟢 Green Flag"},
				}
			}

			raceLive, err2 := sService.GetNextRace(ctx)
			if err2 != nil {
				raceLive = models.Race{Name: "TBD", Date: ""}
			}

			return leaderboard, eventsLive, raceLive, nil
		}

		m := ui.NewLiveModel(fetchFn, 15*time.Second, leaderboard, events, race)
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
