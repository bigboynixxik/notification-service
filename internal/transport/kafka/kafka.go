package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"notification-service/internal/clients"
	"notification-service/internal/models"
	"notification-service/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	auth   *clients.AuthClient
}

func NewConsumer(brokers []string, topic, groupID string, auth *clients.AuthClient) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MaxBytes: 10e6, // 10MB
		}),
		auth: auth,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	l := logger.FromContext(ctx)
	l.Info("kafka consumer started", slog.String("topic", c.reader.Config().Topic))

	for {

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("kafka fetch message: %w", err)
		}

		var event models.NotificationEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			l.Error("failed to unmarshal kafka message", slog.String("error", err.Error()))
			continue
		}

		resp, err := c.auth.GetUserChatID(ctx, event.UserID)
		if err != nil {
			l.Error("failed to get user chat_id",
				slog.String("user_id", event.UserID),
				slog.String("error", err.Error()))
			continue
		}

		l.Info("ready to send notification",
			slog.Int64("chat_id", resp.GetChatId()),
			slog.String("msg", event.Message))

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			l.Error("failed to commit message", slog.String("error", err.Error()))
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
