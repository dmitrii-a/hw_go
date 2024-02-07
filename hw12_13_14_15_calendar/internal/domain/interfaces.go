package domain

import "time"

// EventRepository is an interface for event repository.
type EventRepository interface {
	// Add adds a new event.
	Add(event *Event) error

	// Update updates an existing event.
	Update(event *Event) error

	// Delete removes an event by ID.
	Delete(eventID string) error

	// Get gets an event by ID.
	Get(eventID string) (*Event, error)

	// ListForPeriod get a list of events for a period.
	ListForPeriod(startTime, endTime time.Time) ([]*Event, error)
}
