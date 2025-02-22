package business

import "context"

// Drive adapters
type (
	SendNotificationCase interface {
		SendNotification(context.Context, NotificationRequest) error
	}

	LogNotificationCase interface {
		LogNotification(context.Context, NotificationRequest) error
	}
)

// Driven adapters
type (
	UserStore interface {
		GetUser(context.Context, uint64) (User, error)
	}

	NotificationStore interface {
		SendNotification(context.Context, Notification) error
		SaveNotificationRequest(context.Context, NotificationRequest) error
	}
)
