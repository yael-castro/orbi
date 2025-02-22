//go:build consumer

package container

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/yael-castro/orbi/a/pkg/env"
	"github.com/yael-castro/orbi/b/internal/app/business"
	"github.com/yael-castro/orbi/b/internal/app/handler/consumer"
	"github.com/yael-castro/orbi/b/internal/app/repository/postgres"
	"log"
	"os"
	"strings"
)

func New() Container {
	return &kafkaConsumer{}
}

type kafkaConsumer struct {
	container
	consumer *kafka.Consumer
}

func (k *kafkaConsumer) Inject(ctx context.Context, a any) error {
	switch a := a.(type) {
	case *func(context.Context) error:
		return k.injectConsumer(ctx, a)
	}

	return k.container.Inject(ctx, a)
}

func (k *kafkaConsumer) injectConsumer(ctx context.Context, cmd *func(context.Context) error) (err error) {
	// Logging
	errLogger := log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	infoLogger := log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile)

	// Environment variables
	kafkaServers, err := env.Get("KAFKA_SERVERS")
	if err != nil {
		return
	}

	kafkaTopics, err := env.Get("KAFKA_TOPICS")
	if err != nil {
		return
	}

	kafkaGroup, err := env.Get("KAFKA_GROUP")
	if err != nil {
		return
	}

	// External dependencies
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaServers,
		"enable.auto.commit": "false",
		"group.id":           kafkaGroup,
		"auto.offset.reset":  "earliest",
	})
	if err != nil {
		return err
	}

	// Secondary adapters
	store, err := postgres.NewNotificationStore(infoLogger)
	if err != nil {
		return err
	}

	// Business logic
	cases, err := business.NewLogNotificationCase(store)
	if err != nil {
		return err
	}

	err = c.SubscribeTopics(strings.Split(kafkaTopics, ","), nil)
	if err != nil {
		return err
	}

	// Primary adapters
	*cmd = consumer.NotificationMessages(c, cases, errLogger)
	return
}
