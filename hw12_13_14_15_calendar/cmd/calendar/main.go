package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/http/fiber"
)

func main() {
	log := common.Logger
	server := fiber.NewServer()
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(
			context.Background(), time.Duration(common.Config.Server.ShutdownTimeout)*time.Second,
		)
		defer cancel()
		if err := server.Stop(ctx); common.IsErr(err) {
			log.Error().Msg("failed to stop http server: " + err.Error())
		}
	}()

	log.Info().Msg("calendar service is running...")

	if err := server.Start(ctx); common.IsErr(err) {
		log.Error().Msg("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
