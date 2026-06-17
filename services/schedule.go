package services

import (
	"context"
	"encoding/json"
	"fmt"

	"pitwall/internal/api"
	"pitwall/internal/models"
)

// ScheduleService fetches session/schedule information and converts into domain models.
type ScheduleService struct {
	client *api.Client
}

// NewScheduleService creates a new ScheduleService.
func NewScheduleService(c *api.Client) *ScheduleService {
	return &ScheduleService{client: c}
}

// GetNextRace returns basic race info. Attempts several payload shapes for robustness.
func (s *ScheduleService) GetNextRace(ctx context.Context) (models.Race, error) {
	b, err := s.client.GetCurrentSession(ctx)
	if err != nil {
		return models.Race{}, fmt.Errorf("fetching current session: %w", err)
	}

	// Try common wrapper: { sessions: [ { name, date } ] }
	var wrapped struct {
		Sessions []struct {
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"sessions"`
	}
	if err := json.Unmarshal(b, &wrapped); err == nil && len(wrapped.Sessions) > 0 {
		s0 := wrapped.Sessions[0]
		return models.Race{Name: s0.Name, Date: s0.Date}, nil
	}

	// Fallback: try array of sessions directly
	var arr []struct{
		Name string `json:"name"`
		Date string `json:"date"`
	}
	if err := json.Unmarshal(b, &arr); err == nil && len(arr) > 0 {
		return models.Race{Name: arr[0].Name, Date: arr[0].Date}, nil
	}

	// Last resort: attempt to parse a single object with name/date
	var single struct{
		Name string `json:"name"`
		Date string `json:"date"`
	}
	if err := json.Unmarshal(b, &single); err == nil && (single.Name != "" || single.Date != "") {
		return models.Race{Name: single.Name, Date: single.Date}, nil
	}

	return models.Race{}, fmt.Errorf("unrecognized session payload: %s", string(b))
}
