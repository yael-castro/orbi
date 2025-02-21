package rpc

import (
	"github.com/yael-castro/orbi/b/internal/app/business"
	"github.com/yael-castro/orbi/b/pkg/pb"
)

func BusinessNotification(request *pb.SendNotificationRequest) *business.Notification {
	return &business.Notification{
		UserID:         request.UserId,
		IdempotencyKey: request.IdempotentKey,
	}
}
