package main

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/marcusadriano/tgbot-stt/internal"
	"github.com/marcusadriano/tgbot-stt/pkg/audioconverter"
	"github.com/marcusadriano/tgbot-stt/pkg/fileserver"
	"github.com/marcusadriano/tgbot-stt/pkg/transcript"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {

	time.Local, _ = time.LoadLocation("America/Sao_Paulo")

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(os.Stdout).With().
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

	chatGptApiKey := os.Getenv("WHISPER_GPT_API_KEY")
	tgBotToken := os.Getenv("TG_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Panic(err)
	}

	diskFileServer := fileserver.NewDiskFileserver(&logger, os.TempDir())
	audioConverter := audioconverter.NewFfmpeg(diskFileServer)
	transcriptor := transcript.NewWhisperGptTranscriptor(&logger, chatGptApiKey)

	bot.Debug = os.Getenv("TG_BOT_DEBUG_MODE") == "true"
	logger.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	textHandler := internal.NewTextHandler()
	voiceHandler := internal.NewVoiceHandler(&logger, audioConverter, transcriptor)
	service := internal.NewBotService(&logger, bot, textHandler, voiceHandler)
	service.Start()
}

type tgbotLogger struct {
	logger *zerolog.Logger
}

func (t *tgbotLogger) Printf(format string, v ...interface{}) {
	t.logger.Info().Msgf(format, v...)
}

func (t *tgbotLogger) Println(v ...interface{}) {
	t.logger.Info().Msgf("%s", v...)
}
