package internal

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/marcusadriano/tgbot-stt/internal/logger"
	"github.com/marcusadriano/tgbot-stt/pkg/audioconverter"
	"github.com/marcusadriano/tgbot-stt/pkg/telegram"
	"github.com/marcusadriano/tgbot-stt/pkg/transcript"
)

type voiceHandler struct {
	telegram.Handler
	audioConverter audioconverter.AudioConverter
	transcript     transcript.Transcriptor
}

func NewVoiceHandler(audioConverter audioconverter.AudioConverter, stt transcript.Transcriptor) telegram.Handler {
	vh := voiceHandler{
		audioConverter: audioConverter,
		transcript:     stt,
	}

	return telegram.Handler{
		Handler:   vh.Handle,
		CanHandle: vh.CanHandle,
	}
}

func (v voiceHandler) CanHandle(update tgbotapi.Update) bool {
	return telegram.IsVoice(update)
}

func (v voiceHandler) Handle(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := update.Message

	file, fileData := telegram.DownloadFile(bot, message.Voice.FileID)
	if fileData == nil {
		handleError(bot, update, fmt.Errorf("error to download file"))
		return
	}
	result, err := v.audioConverter.ToMp3(ctx, fileData, file.FilePath)
	if err != nil {
		logger.Log(ctx).Error().Err(err).Msg("Error to convert file")
		handleError(bot, update, err)
		return
	}

	transcription, err := v.transcript.Transcript(ctx, result.Data, "audio.mp3")
	if err != nil {
		logger.Log(ctx).Error().Err(err).Msg("Error to transcript file")
		handleError(bot, update, err)
		return
	}

	transcriptionTextForLog := transcription.Text
	if len(transcriptionTextForLog) > 20 {
		transcriptionTextForLog = transcriptionTextForLog[:20] + "..."
	}

	logger.Log(ctx).Info().
		Str("original_file_path", file.FilePath).
		Msgf("Transcription: %s", transcriptionTextForLog)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, transcription.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)
}
