package application

import (
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
)

type EventService struct {
	repository domain.EventRepository
}

// NewEventService returns a new instance of the event service.
func NewEventService(repository domain.EventRepository) *EventService {
	return &EventService{repository: repository}
}

// Get returns an event by its id.
func (s *EventService) Get(id string) (*domain.Event, error) {
	return s.repository.Get(id)
}

// Create creates a new event.
func (s *EventService) Create(event *domain.Event) error {
	return s.repository.Add(event)
}

// Update updates an existing event.
func (s *EventService) Update(event *domain.Event) error {
	return s.repository.Update(event)
}

// Delete removes an event by ID.
func (s *EventService) Delete(id string) error {
	return s.repository.Delete(id)
}

// ListForPeriod returns a list of events for a period.
func (s *EventService) ListForPeriod(startTime, endTime time.Time) ([]*domain.Event, error) {
	return s.repository.ListForPeriod(startTime, endTime)
}
