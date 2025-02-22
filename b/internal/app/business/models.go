package business

import "fmt"

type NotificationRequest struct {
	UserID         uint64
	IdempotencyKey string
}

func (n NotificationRequest) Validate() error {
	if n.UserID == 0 {
		return fmt.Errorf("%w: missing user id", ErrInvalidNotification)
	}

	return nil
}

type Notification struct {
	Email string
}

type User struct {
	ID    int64
	Name  string
	Email string
	Age   uint8
}

func (u User) Notification() Notification {
	return Notification{
		Email: u.Email,
	}
}
