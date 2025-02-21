package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/yael-castro/orbi/a/internal/app/business"
	"github.com/yael-castro/orbi/a/pkg/jsont"
	"github.com/yael-castro/orbi/b/pkg/pb"
	"google.golang.org/grpc"
	"log"
)

func NewMessageSender(conn *grpc.ClientConn, info *log.Logger) (business.MessageSender, error) {
	if conn == nil {
		return nil, errors.New("conn is nil")
	}

	client := pb.NewNotificationServiceClient(conn)

	return messageSender{
		client: client,
		info:   info,
	}, nil
}

type messageSender struct {
	client pb.NotificationServiceClient
	info   *log.Logger
}

func (m messageSender) SendMessage(ctx context.Context, message *business.Message) error {
	user := jsont.User{}

	err := json.Unmarshal(message.Value, &user)
	if err != nil {
		return err
	}

	idempotencyKey := message.IdempotencyKey

	_, err = m.client.SendNotification(ctx, &pb.SendNotificationRequest{
		UserId:        uint64(user.ID),
		IdempotentKey: string(idempotencyKey),
	})
	if err != nil {
		return err
	}

	m.info.Printf("RPC: message id %d\n", message.ID)
	return nil
}
