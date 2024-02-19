package domain

import (
	"context"
	"io"
	"time"
)

// EventRepository is an interface for event repository.
type EventRepository interface {
	// Add adds a new event.
	Add(event *Event) error

	// Update updates an existing event.
	Update(event *Event) error

	// Delete removes an event by ID.
	Delete(eventID string) error

	// DeleteEventBeforeDate removes an event before date.
	DeleteEventBeforeDate(date time.Time) error

	// Get gets an event by ID.
	Get(eventID string) (*Event, error)

	// GetEventsByPeriod get a list of events for a period.
	GetEventsByPeriod(startTime, endTime time.Time) ([]*Event, error)

	// GetEventsByNotifyTime gets a list of events by notify time.
	GetEventsByNotifyTime(startTime, endTime time.Time) ([]*Event, error)
}

type EventConsumer interface {
	io.Closer
	Consume(name string) (<-chan []byte, error)
}

type EventProducer interface {
	io.Closer
	Publish(ctx context.Context, queueName string, data []byte) error
}
