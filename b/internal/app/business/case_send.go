package business

import (
	"context"
	"fmt"
)

func NewSendNotificationCase(notificationStore NotificationStore, userStore UserStore) (SendNotificationCase, error) {
	if notificationStore == nil || userStore == nil {
		return nil, fmt.Errorf("nil %T", notificationStore)
	}

	return sendNotificationCase{
		userStore:         userStore,
		notificationStore: notificationStore,
	}, nil
}

type sendNotificationCase struct {
	userStore         UserStore
	notificationStore NotificationStore
}

func (n sendNotificationCase) SendNotification(ctx context.Context, notification NotificationRequest) error {
	err := notification.Validate()
	if err != nil {
		return err
	}

	user, err := n.userStore.GetUser(ctx, notification.UserID)
	if err != nil {
		return err
	}

	return n.notificationStore.SendNotification(ctx, user.Notification())
}
