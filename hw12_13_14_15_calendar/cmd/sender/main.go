package main

import (
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/application"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/infrastructure/event"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/infrastructure/repository"
)

func main() {
	common.Config.SetConfigFileSettings(common.GetConfigPathFromArg())
	ctx, cancel := common.GetNotifyCancelCtx()
	defer cancel()
	go func() {
		application.NewEventSenderService(
			repository.GetEventRepository(), event.NewRabbitClient(),
		).Consume(ctx)
	}()
	<-ctx.Done()
}
