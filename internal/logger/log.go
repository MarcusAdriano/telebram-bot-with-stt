package logger

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	RequestIDKey = "request-id"
	ChatIDKey    = "chat-id"
)

var (
	logger *zerolog.Logger
)

func init() {

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log := zerolog.New(os.Stdout).With().
		Caller().
		Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	logger = &log
	zerolog.DefaultContextLogger = nil
}

func Log(ctx context.Context) *zerolog.Logger {

	if ctx == nil {
		return logger
	}

	if ctxLogger := zerolog.Ctx(ctx); ctxLogger != nil {
		return ctxLogger
	}

	return logger
}

func Context(ctx context.Context, chatId int64) context.Context {

	reqId := uuid.New().String()

	newLogger := logger.With().
		Int64(ChatIDKey, chatId).
		Str(RequestIDKey, reqId).
		Logger()

	return newLogger.WithContext(ctx)
}
