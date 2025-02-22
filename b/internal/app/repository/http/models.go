package http

import (
	"github.com/yael-castro/orbi/a/pkg/userapi"
	"github.com/yael-castro/orbi/b/internal/app/business"
)

// User alias for userapi.User
type User = userapi.User

func ToBusiness(u User) business.User {
	return business.User(u)
}
