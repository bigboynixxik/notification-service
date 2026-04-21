package tg

import (
	"context"
	"fmt"
	"log/slog"
	"notification-service/internal/clients"
	"notification-service/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	BotAPI *tgbotapi.BotAPI
	Client *clients.AuthClient
}

func NewBot(token string, client *clients.AuthClient) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("tg.NewBot: %w", err)
	}
	slog.Info("tg.NewBot created bot")
	return &Bot{
		BotAPI: bot,
		Client: client,
	}, nil
}

func (b *Bot) Start(ctx context.Context) {
	l := logger.FromContext(ctx).With("Component", "Telegram bot")
	l.Info("Starting Telegram bot",
		slog.String("account", b.BotAPI.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.BotAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msgCtx := logger.WithContext(ctx, l.With("chat_id", update.Message.Chat.ID,
			"username", update.Message.From.UserName))

		logMsg := logger.FromContext(msgCtx)

		logMsg.Debug("Received message",
			slog.String("text", update.Message.Text))
		if update.Message.IsCommand() {
			b.HandleCommand(msgCtx, update.Message)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не понял. Попробуй команду /start или /help")

		_, err := b.BotAPI.Send(msg)
		if err != nil {
			logMsg.Error("Failed to send message",
				slog.String("err", err.Error()))
		}

	}
}

func (b *Bot) HandleCommand(ctx context.Context, message *tgbotapi.Message) {
	l := logger.FromContext(ctx)
	var responseText string

	switch message.Command() {
	case "start":
		l.Info("User started bot")
		args := message.CommandArguments()
		if len(args) == 0 {
			responseText = "Привет! Я бот проекта Eventify.\nВаш Chat ID: " + fmt.Sprint(message.Chat.ID) + "\nНачинаю привязку токена"
		} else {
			resp, err := b.Client.BindTelegram(ctx, args, message.Chat.ID)
			if err != nil || !resp.GetSuccess() {
				l.Error("tg.Client.BindTelegram:",
					slog.String("error", err.Error()))
				responseText = "Не получилось привязать ваш токен. Попробуйте ещё раз"
			} else {
				responseText = "Ваш аккаунт привязан успешно!"
			}

		}
	case "help":
		responseText = "Я могу отправлять вам уведомления о статусе заказов. Ожидайте новых сообщений!"
	default:
		responseText = "Неизвестная команда. Попробуйте /start или /help"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, responseText)
	_, err := b.BotAPI.Send(msg)
	if err != nil {
		l.Error("tg.HandleCommand:",
			slog.String("error", err.Error()))
	}
}

func (b *Bot) SendMessage(ctx context.Context, chatID int64, message string) error {
	l := logger.FromContext(ctx)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeHTML

	_, err := b.BotAPI.Send(msg)
	if err != nil {
		l.Error("tg.sendMessage, Failed to send message", slog.String("error", err.Error()))
		return fmt.Errorf("tg.sendMessage, Failed to send message: %w", err)
	}
	l.Info("Successfully sent message", slog.Int64("chat_id", chatID))
	return nil
}
