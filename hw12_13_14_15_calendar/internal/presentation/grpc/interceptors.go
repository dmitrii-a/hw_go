package grpc

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func loggingRequestUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	mD, exist := metadata.FromIncomingContext(ctx)
	var userAgent, ip string
	if exist {
		if len(mD["x-forwarded-for"]) > 0 {
			ip = mD["x-forwarded-for"][0]
		}
		if len(mD["user-agent"]) > 0 {
			userAgent = mD["user-agent"][0]
		}
	}
	common.Logger.Info().Msgf(
		"%s [%v] %v %v %v %v \n",
		start.Format(time.RFC3339),
		status.Code(err),
		time.Since(start),
		ip,
		userAgent,
		info.FullMethod,
	)
	return resp, err
}

func recoveryInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = status.Error(codes.Internal, "critical error on server")
			common.Logger.Error().Msgf("panic: %v\n%s", r, debug.Stack())
		}
	}()
	resp, err = handler(ctx, req)
	return resp, err
}
