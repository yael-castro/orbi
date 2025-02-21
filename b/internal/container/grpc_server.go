package container

import (
	"context"
	"github.com/yael-castro/orbi/b/internal/app/business"
	"github.com/yael-castro/orbi/b/internal/app/handler/rpc"
	"github.com/yael-castro/orbi/b/internal/app/repository/postgres"
	"github.com/yael-castro/orbi/b/pkg/interceptor"
	"github.com/yael-castro/orbi/b/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
)

func New() Container {
	return new(containerGRPC)
}

type containerGRPC struct {
	container
}

func (c *containerGRPC) Inject(ctx context.Context, a any) error {
	switch a := a.(type) {
	case **grpc.Server:
		return c.injectServer(ctx, a)
	default:
	}

	return c.container.Inject(ctx, a)
}

func (c *containerGRPC) injectServer(_ context.Context, a **grpc.Server) (err error) {
	info := log.New(os.Stdout, "[INFO]", log.LstdFlags)

	store, err := postgres.NewNotificationStore(info)
	if err != nil {
		return err
	}

	greeter, err := business.NewNotificationCases(store)
	if err != nil {
		return
	}

	service, err := rpc.NewNotificationServiceServer(greeter)
	if err != nil {
		return
	}

	*a = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.Logger(info),
			rpc.ErrorHandling(),
		),
	)

	pb.RegisterNotificationServiceServer(*a, service)
	return
}
