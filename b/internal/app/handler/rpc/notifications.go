package rpc

import (
	"context"
	"fmt"
	"github.com/yael-castro/orbi/b/internal/app/business"
	"github.com/yael-castro/orbi/b/pkg/pb"
)

func NewNotificationServiceServer(notifier business.NotificationCases) (pb.NotificationServiceServer, error) {
	if notifier == nil {
		return nil, fmt.Errorf("%T is nil", notifier)
	}

	return notificationService{
		notifier: notifier,
	}, nil
}

type notificationService struct {
	pb.UnimplementedNotificationServiceServer
	notifier business.NotificationCases
}

func (s notificationService) SendNotification(ctx context.Context, request *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	err := s.notifier.SendNotification(ctx, *BusinessNotification(request))
	if err != nil {
		return nil, err
	}

	response := &pb.SendNotificationResponse{
		Message: new(string),
	}

	*response.Message = "OK"

	return response, nil
}
