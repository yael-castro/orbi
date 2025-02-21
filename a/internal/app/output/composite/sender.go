package composite

import (
	"context"
	"fmt"
	"github.com/yael-castro/orbi/a/internal/app/business"
)

func NewMessageSender(options ...SenderOption) business.MessageSender {
	sender := MessageSender{}

	for _, option := range options {
		option(sender)
	}

	return sender
}

type MessageSender map[string]business.MessageSender

func (m MessageSender) SendMessage(ctx context.Context, message *business.Message) error {
	sender, exists := m[message.Topic]
	if !exists {
		return fmt.Errorf("message sender for topic %q not found", message.Topic)
	}

	return sender.SendMessage(ctx, message)
}

type SenderOption func(MessageSender)

func WithTopic(topic string, sender business.MessageSender) SenderOption {
	return func(message MessageSender) {
		message[topic] = sender
	}
}
