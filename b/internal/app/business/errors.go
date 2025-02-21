package business

import "strconv"

const (
	_ Error = iota
	ErrInvalidNotification
)

type Error uint16

func (e Error) Error() string {
	const errorPrefix = "N"
	return errorPrefix + strconv.Itoa(int(e))
}
