package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/yael-castro/orbi/a/pkg/jsont"
	"github.com/yael-castro/orbi/b/internal/app/business"
	"log"
)

func NotificationMessages(consumer *kafka.Consumer, cases business.NotificationCases, errLogger *log.Logger) func(context.Context) error {
	return func(ctx context.Context) error {
		defer func() {
			_ = consumer.Close()
		}()

		for {
			const pollTimeoutMs = 1_000
			const idempotencyHeader = "idempotency_key"

			// Polling messages
			message, err := consumer.ReadMessage(pollTimeoutMs)
			if err != nil {
				var kafkaErr kafka.Error

				if errors.As(err, &kafkaErr) && kafkaErr.IsTimeout() {
					continue
				}

				errLogger.Printf("Unknown error: %[1]v (%[1]T)", err)
				continue
			}

			// Decoding idempotency key
			var notification business.Notification

			for _, header := range message.Headers {
				if header.Key == idempotencyHeader {
					notification.IdempotencyKey = string(header.Value)
				}
			}

			// Decoding message
			var user jsont.User

			err = json.Unmarshal(message.Value, &user)
			if err != nil {
				errLogger.Println(err)
				// At this point the messages is corrupted or has an invalid structure TODO: send to a DLQ before commit
				goto commit
			}

			notification.UserID = uint64(user.ID)

			// Calling business logic
			// TODO: Should I avoid transforming user data into notification data?
			err = cases.LogNotification(ctx, notification)
			if err != nil {
				errLogger.Println(err)
				continue
			}

		commit:
			_, err = consumer.CommitMessage(message)
			if err != nil {
				errLogger.Println("Fail commit:", err)
				continue
			}
		}
	}
}
