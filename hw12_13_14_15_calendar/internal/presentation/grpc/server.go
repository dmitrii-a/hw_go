package grpc

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation"
	pb "github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc/api/v1"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/grpc/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	grpcServer *grpc.Server
	restServer *http.Server
}

// NewServer returns a new instance of a server.
func NewServer() presentation.Server {
	return &server{}
}

// Start starts the GRPC server.
func (s *server) Start(ctx context.Context) error {
	common.Logger.Info().Msg("grpc service starting...")
	grpcEndpoint := common.GetServerAddr(
		common.Config.Server.GrpcHost,
		common.Config.Server.GrpcPort,
	)
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		return err
	}
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingRequestUnaryInterceptor,
			recoveryInterceptor,
		),
	)
	eventService := service.NewGrpcEventService()
	pb.RegisterEventServiceV1Server(s.grpcServer, eventService)

	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			common.Logger.Fatal().Msgf("grpc ListenAndServe(): %v", err)
		}
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err := pb.RegisterEventServiceV1HandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		return err
	}
	grpcGWEndpoint := common.GetServerAddr(
		common.Config.Server.GrpcGWHost,
		common.Config.Server.GrpcGWPort,
	)
	s.restServer = &http.Server{
		Addr:              grpcGWEndpoint,
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(common.Config.Server.ReadHeaderTimeout) * time.Second,
		ReadTimeout:       time.Duration(common.Config.Server.ReadTimeout) * time.Second,
	}
	go func() {
		if err := s.restServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			common.Logger.Fatal().Msgf("rest grpc ListenAndServe(): %v", err)
		}
	}()
	common.Logger.Info().Msg("grpc service started")
	<-ctx.Done()
	return nil
}

// Stop stops the GRPC server.
func (s *server) Stop(ctx context.Context) error {
	common.Logger.Info().Msg("grpc service is stopping...")
	s.grpcServer.Stop()
	if err := s.restServer.Shutdown(ctx); common.IsErr(err) {
		return err
	}
	return nil
}
