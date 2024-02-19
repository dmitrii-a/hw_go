package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/http/fiber"
)

func main() {
	var path string
	flag.StringVar(&path, "config", "", "Path to configuration file")
	flag.Parse()
	common.Config.SetConfigFileSettings(path)
	log := common.Logger

	server := fiber.NewServer()
	grpcServer := grpc.NewServer()

	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)
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

	log.Info().Msg("calendar service is starting...")

	go func() {
		if err := grpcServer.Start(ctx); common.IsErr(err) {
			log.Error().Msg("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()
	go func() {
		if err := server.Start(ctx); common.IsErr(err) {
			log.Error().Msg("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	<-ctx.Done()
}
