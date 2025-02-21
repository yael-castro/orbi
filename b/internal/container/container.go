package container

import (
	"context"
	"fmt"
)

type Container interface {
	Inject(context.Context, any) error
	Close(context.Context) error
}

type container struct{}

func (c *container) Inject(_ context.Context, a any) error {
	return fmt.Errorf("unsupported type '%T'", a)
}

func (*container) Close(context.Context) error {
	return nil
}
