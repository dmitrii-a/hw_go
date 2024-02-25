package application

import (
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	"github.com/google/uuid"
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
	if err := s.validateID(id); err != nil {
		return nil, err
	}
	return s.repository.Get(id)
}

// Create creates a new event.
func (s *EventService) Create(event *domain.Event) error {
	event.ID = event.NewUUID()
	if err := event.Validate(); common.IsErr(err) {
		return err
	}
	return s.repository.Add(event)
}

// Update updates an existing event.
func (s *EventService) Update(event *domain.Event) error {
	if err := event.Validate(); common.IsErr(err) {
		return err
	}
	return s.repository.Update(event)
}

// Delete removes an event by ID.
func (s *EventService) Delete(id string) error {
	if err := s.validateID(id); err != nil {
		return err
	}
	return s.repository.Delete(id)
}

// ListByPeriod returns a list of events for a period.
func (s *EventService) ListByPeriod(startTime, endTime time.Time) ([]*domain.Event, error) {
	return s.repository.GetEventsByPeriod(startTime, endTime)
}

func (s *EventService) validateID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return domain.ErrUUID
	}
	return nil
}
