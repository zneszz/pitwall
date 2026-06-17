package services

import (
	"encoding/json"
	"fmt"

	"pitwall/internal/models"
)

// parseEvents attempts known JSON shapes and returns domain events or an error.
func parseEvents(b []byte) ([]models.Event, error) {
	// Shape 1: { messages: [ { message: "..." } ] }
	var w struct{
		Messages []struct{ Message string `json:"message"` } `json:"messages"`
	}
	if err := json.Unmarshal(b, &w); err == nil && len(w.Messages) > 0 {
		out := make([]models.Event, 0, len(w.Messages))
		for _, m := range w.Messages {
			out = append(out, models.Event{Message: m.Message})
		}
		return out, nil
	}

	// Shape 2: [ { message: "..." }, ... ]
	var arr []struct{ Message string `json:"message"` }
	if err := json.Unmarshal(b, &arr); err == nil && len(arr) > 0 {
		out := make([]models.Event, 0, len(arr))
		for _, m := range arr {
			out = append(out, models.Event{Message: m.Message})
		}
		return out, nil
	}

	// Shape 3: [ "msg1", "msg2" ]
	var sarr []string
	if err := json.Unmarshal(b, &sarr); err == nil && len(sarr) > 0 {
		out := make([]models.Event, 0, len(sarr))
		for _, s := range sarr {
			out = append(out, models.Event{Message: s})
		}
		return out, nil
	}

	return nil, fmt.Errorf("unrecognized events payload")
}

// parseRace attempts known JSON shapes and returns a domain Race or an error.
func parseRace(b []byte) (models.Race, error) {
	// Shape 1: { sessions: [ { name, date } ] }
	var w struct{
		Sessions []struct{
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"sessions"`
	}
	if err := json.Unmarshal(b, &w); err == nil && len(w.Sessions) > 0 {
		s := w.Sessions[0]
		return models.Race{Name: s.Name, Date: s.Date}, nil
	}

	// Shape 2: [ { name, date }, ... ]
	var arr []struct{
		Name string `json:"name"`
		Date string `json:"date"`
	}
	if err := json.Unmarshal(b, &arr); err == nil && len(arr) > 0 {
		return models.Race{Name: arr[0].Name, Date: arr[0].Date}, nil
	}

	// Shape 3: single object { name, date }
	var s struct{
		Name string `json:"name"`
		Date string `json:"date"`
	}
	if err := json.Unmarshal(b, &s); err == nil && (s.Name != "" || s.Date != "") {
		return models.Race{Name: s.Name, Date: s.Date}, nil
	}

	return models.Race{}, fmt.Errorf("unrecognized race payload")
}
