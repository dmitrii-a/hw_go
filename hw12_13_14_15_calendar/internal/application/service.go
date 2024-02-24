package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
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
	return s.repository.GetEventsByPeriod(startTime, endTime)
}

type EventSchedulerService struct {
	repository domain.EventRepository
	producer   domain.EventProducer
	consumer   domain.EventConsumer
}

// NewEventSchedulerService returns a new instance of the event scheduler service.
func NewEventSchedulerService(
	repository domain.EventRepository,
	producer domain.EventProducer,
) *EventSchedulerService {
	return &EventSchedulerService{repository: repository, producer: producer}
}

// NewEventSenderService returns a new instance of the event sender service.
func NewEventSenderService(
	repository domain.EventRepository,
	consumer domain.EventConsumer,
) *EventSchedulerService {
	return &EventSchedulerService{repository: repository, consumer: consumer}
}

func (s *EventSchedulerService) cleanEvents() {
	common.Logger.Info().Msg("Start clean events")
	t := time.Now().Add(-time.Duration(common.Config.Scheduler.EventLifetime) * time.Second)
	err := s.repository.DeleteEventBeforeDate(t)
	if common.IsErr(err) {
		common.Logger.Error().Msgf("failed to clean events: %v", err)
	}
}

// Schedule schedules events.
func (s *EventSchedulerService) Schedule(ctx context.Context) {
	common.Logger.Info().Msg("start scheduler")
	periodTime := time.Duration(common.Config.Scheduler.PublishPeriodTime) * time.Second
	ticker := time.NewTicker(periodTime)
	startDate := time.Now().UTC().Add(-periodTime).Round(time.Second)
	for {
		select {
		case <-ctx.Done():
			err := s.producer.Close()
			if common.IsErr(err) {
				common.Logger.Error().Msgf("failed to close producer: %v", err)
			}
		case <-ticker.C:
			endDate := time.Now().UTC().Round(time.Second)
			events, err := s.repository.GetEventsByNotifyTime(startDate, endDate)
			if common.IsErr(err) {
				common.Logger.Error().Msgf("failed to handle events: %v", err)
			}
			notifications := make([]*domain.Notification, len(events))
			for i, event := range events {
				notification := domain.Notification{
					EventID:    event.ID,
					EventTitle: event.Title,
					EventDate:  event.StartTime,
					UserToSend: event.UserID,
				}
				notifications[i] = &notification
			}
			data, err := json.Marshal(notifications)
			if common.IsErr(err) {
				common.Logger.Error().Msgf("failed to marshal notification: %v", err)
				continue
			}
			common.Logger.Info().Msg("start publishing notifications")
			err = s.producer.Publish(ctx, "events", data)
			if common.IsErr(err) {
				common.Logger.Error().Msgf("failed to publish notification: %v", err)
				continue
			}
			startDate = endDate
			s.cleanEvents()
		}
	}
}

// Consume consumes events.
func (s *EventSchedulerService) Consume(ctx context.Context) {
	common.Logger.Info().Msg("start consume")
	consumer, err := s.consumer.Consume("events")
	if common.IsErr(err) {
		common.Logger.Error().Msgf("failed to consume: %v", err)
	}
	for {
		select {
		case <-ctx.Done():
			err := s.consumer.Close()
			if common.IsErr(err) {
				common.Logger.Error().Msgf("failed to close consumer: %v", err)
			}
		case data := <-consumer:
			var notifications []*domain.Notification
			err := json.Unmarshal(data, &notifications)
			if common.IsErr(err) {
				common.Logger.Error().Msgf("failed to unmarshal notification: %v", err)
				continue
			}
			for _, notification := range notifications {
				common.Logger.Info().Msgf("send notification: %v", notification)
			}
		}
	}
}
