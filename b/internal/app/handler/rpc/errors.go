package rpc

import (
	"context"
	"errors"
	"github.com/yael-castro/orbi/b/internal/app/business"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandling() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)

		// Checking if there are a business error
		var businessErr business.Error

		if !errors.As(err, &businessErr) {
			return resp, err
		}

		//goland:noinspection ALL
		switch businessErr {
		case business.ErrInvalidNotification:
			return nil, status.Errorf(codes.InvalidArgument, "%v", err)
		}

		return resp, err
	}
}
