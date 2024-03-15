package application

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
)

type EventSchedulerProcessor struct {
	repository domain.EventRepository
	producer   domain.EventProducer
	consumer   domain.EventConsumer
}

const (
	EventQueueName       = "events"
	EventResultQueueName = "events_result"
)

// NewEventSchedulerProcessor returns a new instance of the event scheduler service.
func NewEventSchedulerProcessor(
	repository domain.EventRepository,
	producer domain.EventProducer,
) *EventSchedulerProcessor {
	return &EventSchedulerProcessor{repository: repository, producer: producer}
}

// NewEventSenderProcessor returns a new instance of the event sender service.
func NewEventSenderProcessor(
	repository domain.EventRepository,
	consumer domain.EventConsumer,
	producer domain.EventProducer,
) *EventSchedulerProcessor {
	return &EventSchedulerProcessor{repository: repository, consumer: consumer, producer: producer}
}

func (s *EventSchedulerProcessor) cleanEvents() {
	common.Logger.Info().Msg("event cleanup started")
	t := time.Now().Add(-time.Duration(common.Config.Scheduler.EventLifetime) * time.Second)
	err := s.repository.DeleteEventBeforeDate(t)
	if common.IsErr(err) {
		common.Logger.Error().Msgf("failed to clean events: %v", err)
	}
	common.Logger.Info().Msg("event cleanup completed")
}

// Schedule schedules events.
func (s *EventSchedulerProcessor) Schedule(ctx context.Context) {
	common.Logger.Info().Msg("running the scheduler")
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
			common.Logger.Info().Msg("started publishing notifications")
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
func (s *EventSchedulerProcessor) Consume(ctx context.Context) {
	common.Logger.Info().Msg("start consume")
	consumer, err := s.consumer.Consume(EventQueueName)
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
				common.Logger.Info().Msgf("sending notification: %v", notification)
				err = s.producer.Publish(ctx, EventResultQueueName, []byte(notification.EventID))
				if common.IsErr(err) {
					common.Logger.Error().Msgf("failed to publish result in queue: %v", err)
				}
			}
		}
	}
}
