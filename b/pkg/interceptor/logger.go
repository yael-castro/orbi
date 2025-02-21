package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
)

func Logger(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)

		state, _ := status.FromError(err)

		code := state.Code()

		logger.Printf("%s %v %s (%v)\n", info.FullMethod, resp, code, int(code))
		return
	}
}
