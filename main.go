package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/marcusadriano/sound-stt-tgbot/internal/audioconverter"
	"github.com/marcusadriano/sound-stt-tgbot/internal/transcript"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	chatGptApiKey  string
	logger         zerolog.Logger
	audioConverter audioconverter.AudioConverter
	transcriptor   transcript.Transcriptor
)

func main() {

	time.Local, _ = time.LoadLocation("America/Sao_Paulo")

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger = zerolog.New(os.Stdout).With().
		Caller().
		Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	botLogger := &tgbotLogger{logger: &logger}
	tgbotapi.SetLogger(botLogger)

	chatGptApiKey = os.Getenv("WHISPER_GPT_API_KEY")
	tgBotToken := os.Getenv("TG_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Panic(err)
	}

	audioConverter = audioconverter.NewFfmpeg(&logger)
	transcriptor = transcript.NewWhisperGptTranscriptor(&logger, chatGptApiKey)

	bot.Debug = os.Getenv("TG_BOT_DEBUG_MODE") == "true"

	logger.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		go defaultHandle(bot, update)
	}
}

func defaultHandle(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	logger.Info().Msgf("Receive a message from %s", update.Message.From.UserName)

	if update.Message != nil {

		message := update.Message
		if message.Voice != nil {
			go handleAudio(bot, update)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func handleAudio(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := update.Message
	if message.Voice != nil {
		go func() {

			ctx := context.Background()

			file, fileData := downloadTelegramVoice(bot, message, update)
			if fileData == nil {
				return
			}
			result, err := audioConverter.ToMp3(ctx, fileData, file.FilePath)
			if err != nil {
				logger.Error().Err(err).Msg("Error to convert file")
				handleError(bot, update, err)
				return
			}

			transcription, err := transcriptor.Transcript(result.Data, "audio.mp3")
			if err != nil {
				logger.Error().Err(err).Msg("Error to transcript file")
				handleError(bot, update, err)
				return
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, transcription.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}()
	}
}

func downloadTelegramVoice(bot *tgbotapi.BotAPI, message *tgbotapi.Message, update tgbotapi.Update) (tgbotapi.File, []byte) {
	file, err := bot.GetFile(tgbotapi.FileConfig{
		FileID: message.Voice.FileID,
	})

	if err != nil {
		logger.Error().Err(err).Msg("Error to download file")
		handleError(bot, update, err)
		return tgbotapi.File{}, nil
	}

	fileLink := file.Link(bot.Token)
	logger.Info().Msgf("Downloading file %s", fileLink)

	req, _ := http.NewRequest("GET", fileLink, nil)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		logger.Error().Err(err).Msg("Error to download file")
		handleError(bot, update, err)
		return tgbotapi.File{}, nil
	}
	defer res.Body.Close()

	fileData, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Error to download file")
		handleError(bot, update, err)
		return tgbotapi.File{}, nil
	}
	return file, fileData
}

func handleError(bot *tgbotapi.BotAPI, update tgbotapi.Update, err error) {

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I can't listen you")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}

type tgbotLogger struct {
	logger *zerolog.Logger
}

func (t *tgbotLogger) Printf(format string, v ...interface{}) {
	t.logger.Info().Msgf(format, v)
}

func (t *tgbotLogger) Println(v ...interface{}) {
	t.logger.Info().Msgf("%s", v)
}
