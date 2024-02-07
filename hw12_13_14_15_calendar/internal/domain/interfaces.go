package domain

import "time"

// EventRepository is an interface for event repository.
type EventRepository interface {
	// AddEvent adds a new event.
	AddEvent(event *Event) error

	// UpdateEvent updates an existing event.
	UpdateEvent(event *Event) error

	// DeleteEvent removes an event by ID.
	DeleteEvent(eventID string) error

	// GetEvent get an event by ID.
	GetEvent(eventID string) (*Event, error)

	// ListEventsForPeriod get a list of events for a period.
	ListEventsForPeriod(startTime, endTime time.Time) ([]*Event, error)
}
