package tg

import (
	"context"
	"testing"

	"notification-service/internal/clients"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestBot_HandleCommand_Coverage(t *testing.T) {
	bot := &Bot{
		BotAPI: &tgbotapi.BotAPI{Self: tgbotapi.User{UserName: "test_bot"}},
		Client: &clients.AuthClient{},
	}

	tests := []struct {
		name    string
		command string
		args    string
	}{
		{"Команда start без аргументов", "start", ""},
		{"Команда help", "help", ""},
		{"Неизвестная команда", "unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			message := &tgbotapi.Message{
				Text: "/" + tt.command,
				Entities: []tgbotapi.MessageEntity{
					{Type: "bot_command", Offset: 0, Length: len(tt.command) + 1},
				},
				Chat: &tgbotapi.Chat{ID: 12345},
			}

			bot.HandleCommand(ctx, message)
		})
	}
}
