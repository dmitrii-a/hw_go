package application

import (
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/infrastructure/repository"
)

// EventApplicationService instance of the event service.
var EventApplicationService *EventService

func init() {
	var eventRepository domain.EventRepository
	if common.Config.UseCacheDB {
		eventRepository = repository.NewEventCacheRepository()
	} else {
		eventRepository = repository.NewEventDBRepository()
	}
	EventApplicationService = NewEventService(eventRepository)
}
