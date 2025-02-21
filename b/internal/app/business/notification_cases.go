package business

import (
	"context"
	"fmt"
)

func NewNotificationCases(store NotificationStore) (NotificationCases, error) {
	if store == nil {
		return nil, fmt.Errorf("nil %T", store)
	}

	return notificationCases{
		store: store,
	}, nil
}

type notificationCases struct {
	store NotificationStore
}

func (n notificationCases) SendNotification(ctx context.Context, notification Notification) error {
	err := notification.Validate()
	if err != nil {
		return err
	}

	// TODO: Retrieve user data

	return n.store.SendNotification(ctx, notification)
}

func (n notificationCases) LogNotification(ctx context.Context, notification Notification) error {
	err := notification.Validate()
	if err != nil {
		return err
	}

	return n.store.SaveNotification(ctx, notification)
}
