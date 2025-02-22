package postgres

import (
	"context"
	"errors"
	"github.com/yael-castro/orbi/b/internal/app/business"
	"log"
)

func NewNotificationStore(info *log.Logger) (business.NotificationStore, error) {
	if info == nil {
		return nil, errors.New("logger is required")
	}

	return notificationStore{
		info: info,
	}, nil
}

type notificationStore struct {
	info *log.Logger
}

func (n notificationStore) SendNotification(_ context.Context, notification business.Notification) error {
	n.info.Printf("SendNotification: %+v", notification)
	return nil
}

func (n notificationStore) SaveNotificationRequest(_ context.Context, request business.NotificationRequest) error {
	n.info.Printf("SaveNotificationRequest: %+v\n", request)
	return nil
}
