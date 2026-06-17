package services

import (
	"context"
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
	r, err := parseRace(b)
	if err != nil {
		return models.Race{}, fmt.Errorf("parsing session payload: %w", err)
	}
	return r, nil
}
