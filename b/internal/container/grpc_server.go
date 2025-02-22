//go:build grpc

package container

import (
	"context"
	"github.com/yael-castro/orbi/a/pkg/env"
	"github.com/yael-castro/orbi/a/pkg/userapi"
	"github.com/yael-castro/orbi/b/internal/app/business"
	"github.com/yael-castro/orbi/b/internal/app/handler/rpc"
	"github.com/yael-castro/orbi/b/internal/app/repository/http"
	"github.com/yael-castro/orbi/b/internal/app/repository/postgres"
	"github.com/yael-castro/orbi/b/pkg/interceptor"
	"github.com/yael-castro/orbi/b/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"
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

func (c *containerGRPC) injectServer(ctx context.Context, a **grpc.Server) (err error) {
	// Environment variables
	address, err := env.Get("USERS_API_ADDRESS")
	if err != nil {
		return err
	}

	// Building external dependencies
	api, err := userapi.New(address)
	if err != nil {
		return err
	}

	const pingTimeout = time.Second * 5
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	err = api.Ping(ctx)
	if err != nil {
		return err
	}

	info := log.New(os.Stdout, "[INFO]", log.LstdFlags)

	// Building driven adapters
	notificationStore, err := postgres.NewNotificationStore(info)
	if err != nil {
		return err
	}

	userStore, err := http.NewUserStore(api)
	if err != nil {
		return err
	}

	// Building business logic
	sender, err := business.NewSendNotificationCase(notificationStore, userStore)
	if err != nil {
		return
	}

	// Building drive adapters
	service, err := rpc.NewNotificationServiceServer(sender)
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
