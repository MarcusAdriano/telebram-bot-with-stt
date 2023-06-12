package internal

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/marcusadriano/tgbot-stt/pkg/telegram"
)

func NewTextHandler() telegram.Handler {
	return telegram.Handler{
		Handler:   textHandlerFunc,
		CanHandle: telegram.IsTextMessage,
	}
}

func textHandlerFunc(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	messageReceived := update.Message.Text
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, messageReceived))
}
