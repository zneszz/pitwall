package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"pitwall/internal/api"
	"pitwall/internal/models"
)

// StandingsService fetches driver and team standings.
type StandingsService struct {
	client *api.Client
}

func NewStandingsService(c *api.Client) *StandingsService {
	return &StandingsService{client: c}
}

// GetDriverStandings attempts to parse common shapes for driver standings.
func (s *StandingsService) GetDriverStandings(ctx context.Context) ([]models.DriverStanding, error) {
	b, err := s.client.GetDriverStandings(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching driver standings: %w", err)
	}

	// Attempt shape: { standings: [ { position: 1, driver: { name: "" }, team: { name: "" }, points: "123" } ] }
	var wrapper struct {
		Standings []struct {
			Position int `json:"position"`
			Driver   struct {
				Name string `json:"name"`
			} `json:"driver"`
			Team struct {
				Name string `json:"name"`
			} `json:"team"`
			Points interface{} `json:"points"`
		} `json:"standings"`
	}
	if err := json.Unmarshal(b, &wrapper); err == nil && len(wrapper.Standings) > 0 {
		out := make([]models.DriverStanding, 0, len(wrapper.Standings))
		for _, s0 := range wrapper.Standings {
			pts := parsePoints(s0.Points)
			out = append(out, models.DriverStanding{Position: s0.Position, DriverName: s0.Driver.Name, Team: s0.Team.Name, Points: pts})
		}
		return out, nil
	}

	// Fallback: try array of entries
	var arr []struct {
		Position int         `json:"position"`
		Name     string      `json:"name"`
		Team     string      `json:"team"`
		Points   interface{} `json:"points"`
	}
	if err := json.Unmarshal(b, &arr); err == nil && len(arr) > 0 {
		out := make([]models.DriverStanding, 0, len(arr))
		for _, a := range arr {
			out = append(out, models.DriverStanding{Position: a.Position, DriverName: a.Name, Team: a.Team, Points: parsePoints(a.Points)})
		}
		return out, nil
	}

	return nil, fmt.Errorf("unrecognized driver standings payload: %s", string(b))
}

// GetTeamStandings attempts to parse team standings.
func (s *StandingsService) GetTeamStandings(ctx context.Context) ([]models.TeamStanding, error) {
	b, err := s.client.GetTeamStandings(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching team standings: %w", err)
	}

	var wrapper struct {
		Standings []struct {
			Position int `json:"position"`
			Team     struct {
				Name string `json:"name"`
			} `json:"team"`
			Points interface{} `json:"points"`
		} `json:"standings"`
	}
	if err := json.Unmarshal(b, &wrapper); err == nil && len(wrapper.Standings) > 0 {
		out := make([]models.TeamStanding, 0, len(wrapper.Standings))
		for _, s0 := range wrapper.Standings {
			out = append(out, models.TeamStanding{Position: s0.Position, TeamName: s0.Team.Name, Points: parsePoints(s0.Points)})
		}
		return out, nil
	}

	var arr []struct {
		Position int         `json:"position"`
		Name     string      `json:"name"`
		Points   interface{} `json:"points"`
	}
	if err := json.Unmarshal(b, &arr); err == nil && len(arr) > 0 {
		out := make([]models.TeamStanding, 0, len(arr))
		for _, a := range arr {
			out = append(out, models.TeamStanding{Position: a.Position, TeamName: a.Name, Points: parsePoints(a.Points)})
		}
		return out, nil
	}

	return nil, fmt.Errorf("unrecognized team standings payload: %s", string(b))
}

// parsePoints converts interface{} points into float64; supports string or number
func parsePoints(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case string:
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f
		}
	}
	return 0.0
}
