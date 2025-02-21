package business

import "fmt"

type User struct{}

type Notification struct {
	UserID         uint64
	IdempotencyKey string
}

func (n Notification) Validate() error {
	if n.UserID == 0 {
		return fmt.Errorf("%w: missing user id", ErrInvalidNotification)
	}

	return nil
}
