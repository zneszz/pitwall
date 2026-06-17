package services

import (
	"context"
	"fmt"

	"pitwall/internal/api"
	"pitwall/internal/models"
)

// EventsService fetches and converts race control messages into domain models.
type EventsService struct {
	client *api.Client
}

// NewEventsService creates a new EventsService.
func NewEventsService(c *api.Client) *EventsService {
	return &EventsService{client: c}
}

// GetRecentEvents returns a slice of domain Event objects.
func (s *EventsService) GetRecentEvents(ctx context.Context) ([]models.Event, error) {
	b, err := s.client.GetRaceControlMessages(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching race control messages: %w", err)
	}
	evs, err := parseEvents(b)
	if err != nil {
		return nil, fmt.Errorf("parsing race control messages: %w", err)
	}
	return evs, nil
}
