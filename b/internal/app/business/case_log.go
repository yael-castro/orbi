package business

import (
	"context"
	"errors"
)

func NewLogNotificationCase(store NotificationStore) (LogNotificationCase, error) {
	if store == nil {
		return nil, errors.New("nil notification store")
	}

	return logNotificationCase{
		store: store,
	}, nil
}

type logNotificationCase struct {
	store NotificationStore
}

func (n logNotificationCase) LogNotification(ctx context.Context, notification NotificationRequest) error {
	err := notification.Validate()
	if err != nil {
		return err
	}

	return n.store.SaveNotificationRequest(ctx, notification)
}
