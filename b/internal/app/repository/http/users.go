package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/yael-castro/orbi/a/pkg/userapi"
	"github.com/yael-castro/orbi/b/internal/app/business"
)

func NewUserStore(api userapi.UserAPI) (business.UserStore, error) {
	if api == nil {
		return nil, errors.New("api required")
	}

	return &userStore{api}, nil
}

type userStore struct {
	api userapi.UserAPI
}

func (u userStore) GetUser(ctx context.Context, userID uint64) (business.User, error) {
	apiUser, err := u.api.GetUser(ctx, userID)
	if err != nil {
		// TODO: improve error handling
		return business.User{}, fmt.Errorf("%w: user %d not found", userID, business.ErrResourceNotFound)
	}

	return ToBusiness(apiUser), nil
}
