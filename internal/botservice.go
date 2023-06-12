package internal

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/marcusadriano/tgbot-stt/internal/logger"
	"github.com/marcusadriano/tgbot-stt/pkg/telegram"
)

type BotService struct {
	bot      *tgbotapi.BotAPI
	handlers *telegram.TgBotHandlers
}

func handleError(bot *tgbotapi.BotAPI, update tgbotapi.Update, err error) {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, There is an error to process your request.")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

func NewBotService(bot *tgbotapi.BotAPI, handlers ...telegram.Handler) *BotService {

	botHandlers := telegram.NewTgBotHandlersBuilder().
		Bot(bot).
		AddHandlers(handlers...).
		Build()

	return &BotService{
		bot:      bot,
		handlers: botHandlers,
	}
}

func (b *BotService) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {

		ctx := createRequestContext(update)
		logger.Log(ctx).Info().Msg("Received new update")

		go b.handlers.Handle(ctx, update)
	}
}

func createRequestContext(update tgbotapi.Update) context.Context {
	ctx := context.Background()
	return logger.Context(ctx, update.Message.Chat.ID)
}
