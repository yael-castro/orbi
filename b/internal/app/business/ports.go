package business

import "context"

type (
	NotificationCases interface {
		LogNotification(context.Context, Notification) error
		SendNotification(context.Context, Notification) error
	}
)

type (
	NotificationStore interface {
		SendNotification(context.Context, Notification) error
		SaveNotification(context.Context, Notification) error
	}
)
