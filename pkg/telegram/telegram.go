package telegram

import (
	"context"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update)
type CanHandleFunc func(update tgbotapi.Update) bool

type Handler struct {
	Handler   HandlerFunc
	CanHandle CanHandleFunc
}

type TgBotHandlers struct {
	bot      *tgbotapi.BotAPI
	handlers []Handler
}

type TgBotHandlersBuilder struct {
	bot      *tgbotapi.BotAPI
	handlers []Handler
}

func NewTgBotHandlersBuilder() *TgBotHandlersBuilder {
	return &TgBotHandlersBuilder{
		handlers: make([]Handler, 0),
	}
}

func (b *TgBotHandlersBuilder) AddHandler(handler Handler) *TgBotHandlersBuilder {
	b.handlers = append(b.handlers, handler)
	return b
}

func (b *TgBotHandlersBuilder) AddHandlers(handlers ...Handler) *TgBotHandlersBuilder {
	b.handlers = append(b.handlers, handlers...)
	return b
}

func (b *TgBotHandlersBuilder) Bot(bot *tgbotapi.BotAPI) *TgBotHandlersBuilder {
	b.bot = bot
	return b
}

func (b *TgBotHandlersBuilder) Build() *TgBotHandlers {
	return &TgBotHandlers{
		bot:      b.bot,
		handlers: b.handlers,
	}
}

func DownloadFile(bot *tgbotapi.BotAPI, fileId string) (tgbotapi.File, []byte) {

	file, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: fileId,
	})

	if err != nil {
		return tgbotapi.File{}, nil
	}

	fileLink := file.Link(bot.Token)

	req, _ := http.NewRequest("GET", fileLink, nil)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return tgbotapi.File{}, nil
	}
	defer res.Body.Close()

	fileData, err := io.ReadAll(res.Body)
	if err != nil {
		return tgbotapi.File{}, nil
	}
	return file, fileData
}

func (h *TgBotHandlers) Handle(ctx context.Context, update tgbotapi.Update) {
	for _, handler := range h.handlers {
		if handler.CanHandle(update) {
			handler.Handler(ctx, h.bot, update)
			return
		}
	}
}

func IsTextMessage(update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text != ""
}

func IsVoice(update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Voice != nil
}
