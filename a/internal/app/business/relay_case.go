//go:build relay

package business

import (
	"context"
	"errors"
	"log"
	"time"
)

type MessagesRelayConfig struct {
	Confirmer   MessageDeliveryConfirmer
	Reader      MessagesReader
	Sender      MessageSender
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

func NewMessagesRelay(config MessagesRelayConfig) (MessagesRelay, error) {
	// TODO: avoid a line too long
	if config.Confirmer == nil || config.InfoLogger == nil || config.Reader == nil || config.Sender == nil || config.ErrorLogger == nil {
		return nil, errors.New("missing settings")
	}

	return &messagesRelay{
		confirmer: config.Confirmer,
		reader:    config.Reader,
		sender:    config.Sender,
		info:      config.InfoLogger,
		error:     config.ErrorLogger,
	}, nil
}

type messagesRelay struct {
	confirmer MessageDeliveryConfirmer
	reader    MessagesReader
	sender    MessageSender
	info      *log.Logger
	error     *log.Logger
}

func (m *messagesRelay) RelayMessages(ctx context.Context) (err error) {
	var messages []Message

	for {
		messages, err = m.reader.ReadMessages(ctx)
		if err != nil {
			return
		}

		// Waiting to poll more messages
		if len(messages) == 0 {
			const retryDelay = 100 * time.Millisecond

			select {
			case <-ctx.Done():
				return errors.Join(err, m.reader.Close())
			case <-time.After(retryDelay):
				continue
			}
		}

		// Relaying messages...
		err = m.relayMessages(ctx, messages)
		if err != nil {
			return
		}
	}
}

func (m *messagesRelay) relayMessages(ctx context.Context, messages []Message) error {
	m.info.Printf("Relaying %d messages...", len(messages))

	for i := range messages {
		err := m.sender.SendMessage(ctx, &messages[i])
		if err != nil {
			// TODO: It is better to have duplicates than an undelivered message (a message is acknowledged when it was not sent).
			//if errors.Is(err, ErrMessageDeliveryFailed) {
			//	return err
			//}

			m.error.Printf("Message %d: '%[2]v' %[2]T", messages[i].ID, err)
			continue
		}

		m.info.Printf("Message %d sent", messages[i].ID)

		err = m.confirmer.ConfirmMessageDelivery(ctx, messages[i].ID)
		if err != nil {
			m.error.Printf("Failed to confirm message %d\n", messages[i].ID)
			continue
		}
	}

	return nil
}
