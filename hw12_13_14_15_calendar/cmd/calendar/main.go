package main

import (
	"context"
	"os"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/http/fiber"
)

func main() {
	common.Config.SetConfigFileSettings(common.GetConfigPathFromArg())

	server := fiber.NewServer()
	grpcServer := grpc.NewServer()

	ctx, cancel := common.GetNotifyCancelCtx()
	defer cancel()

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(
			context.Background(), time.Duration(common.Config.Server.ShutdownTimeout)*time.Second,
		)
		defer cancel()
		if err := server.Stop(ctx); common.IsErr(err) {
			common.Logger.Error().Msg("failed to stop http server: " + err.Error())
		}
	}()

	common.Logger.Info().Msg("calendar service is starting...")

	go func() {
		if err := grpcServer.Start(ctx); common.IsErr(err) {
			common.Logger.Error().Msg("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()
	go func() {
		if err := server.Start(ctx); common.IsErr(err) {
			common.Logger.Error().Msg("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()
	<-ctx.Done()
}
