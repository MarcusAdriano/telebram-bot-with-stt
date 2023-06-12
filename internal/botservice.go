package internal

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/marcusadriano/tgbot-stt/pkg/telegram"
	"github.com/rs/zerolog"
)

type BotService struct {
	logger   *zerolog.Logger
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

func NewBotService(logger *zerolog.Logger, bot *tgbotapi.BotAPI, handlers ...telegram.Handler) *BotService {

	botHandlers := telegram.NewTgBotHandlersBuilder().
		Bot(bot).
		AddHandlers(handlers...).
		Build()

	return &BotService{
		logger:   logger,
		bot:      bot,
		handlers: botHandlers,
	}
}

func (b *BotService) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {

		ctx := context.Background()

		b.logger.Info().
			Int64("chat-id", update.Message.Chat.ID).
			Msgf("Message received")

		go b.handlers.Handle(ctx, update)
	}
}
