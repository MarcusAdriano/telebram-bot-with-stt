package main

import (
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/marcusadriano/tgbot-stt/internal"
	"github.com/marcusadriano/tgbot-stt/internal/logger"
	"github.com/marcusadriano/tgbot-stt/pkg/audioconverter"
	"github.com/marcusadriano/tgbot-stt/pkg/fileserver"
	"github.com/marcusadriano/tgbot-stt/pkg/transcript"
)

func main() {

	time.Local, _ = time.LoadLocation("America/Sao_Paulo")

	err := godotenv.Load()
	if err != nil {
		logger.Default().Fatal().Err(err).Msg("Error loading .env file")
		return
	}

	chatGptApiKey := os.Getenv("WHISPER_GPT_API_KEY")
	tgBotToken := os.Getenv("TG_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		logger.Default().Panic().Err(err).Msg("Error creating bot")
	}

	diskFileServer := fileserver.NewDiskFileserver(os.TempDir())
	audioConverter := audioconverter.NewFfmpeg(diskFileServer)
	transcriptor := transcript.NewWhisperGptTranscriptor(chatGptApiKey)

	bot.Debug = os.Getenv("TG_BOT_DEBUG_MODE") == "true"
	logger.Default().Info().Msgf("Authorized on account %s", bot.Self.UserName)

	textHandler := internal.NewTextHandler()
	voiceHandler := internal.NewVoiceHandler(audioConverter, transcriptor)
	service := internal.NewBotService(bot, textHandler, voiceHandler)
	service.Start()
}
