package services

import (
	"context"
	"encoding/json"
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

	// Try to unmarshal into a wrapper { messages: [...] }
	var wrapped struct {
		Messages []struct {
			Message string `json:"message"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(b, &wrapped); err == nil && len(wrapped.Messages) > 0 {
		out := make([]models.Event, 0, len(wrapped.Messages))
		for _, m := range wrapped.Messages {
			out = append(out, models.Event{Message: m.Message})
		}
		return out, nil
	}

	// Fallback: try to unmarshal as an array of objects with message field
	var arr []struct{
		Message string `json:"message"`
	}
	if err := json.Unmarshal(b, &arr); err == nil && len(arr) > 0 {
		out := make([]models.Event, 0, len(arr))
		for _, m := range arr {
			out = append(out, models.Event{Message: m.Message})
		}
		return out, nil
	}

	// Last resort: attempt to parse as array of strings
	var sarr []string
	if err := json.Unmarshal(b, &sarr); err == nil && len(sarr) > 0 {
		out := make([]models.Event, 0, len(sarr))
		for _, m := range sarr {
			out = append(out, models.Event{Message: m})
		}
		return out, nil
	}

	return nil, fmt.Errorf("unrecognized race control payload: %s", string(b))
}
