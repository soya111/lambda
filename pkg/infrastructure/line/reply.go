package line

import (
	"context"
	"fmt"
	"zephyr/pkg/logging"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

// ReplyTextMessages reply text messages to LINE
func (b *Linebot) ReplyTextMessages(ctx context.Context, token string, message string) error {
	logger := logging.LoggerFromContext(ctx)
	if _, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(message)).WithContext(ctx).Do(); err != nil {
		logger.Error("Failed to reply text message to LINE.", zap.Error(err), zap.String("token", token), zap.String("message", message))
		return err
	}
	logger.Info("Successfully replied text message to LINE.", zap.String("token", token), zap.String("message", message))
	return nil
}

// ReplyMessage reply messages to LINE
func (b *Linebot) ReplyMessage(ctx context.Context, token string, messages ...linebot.SendingMessage) error {
	logger := logging.LoggerFromContext(ctx)
	if _, err := b.Client.ReplyMessage(token, messages...).WithContext(ctx).Do(); err != nil {
		logger.Error("Failed to reply messages to LINE.", zap.Error(err), zap.String("token", token), zap.Any("messages", messages))
		return err
	}
	logger.Info("Successfully replied messages to LINE.", zap.String("token", token))
	return nil
}

// ReplyWithError reply message to LINE and return error
func (b *Linebot) ReplyWithError(ctx context.Context, token, replyMessage string, err error) error {
	logger := logging.LoggerFromContext(ctx)
	if err := b.ReplyTextMessages(ctx, token, replyMessage); err != nil {
		logger.Error("Failed to reply with error message to LINE.", zap.Error(err), zap.String("token", token), zap.String("replyMessage", replyMessage))
		return fmt.Errorf("replyWithError: %w", err)
	}
	logger.Info("Successfully replied with error message to LINE.", zap.String("token", token), zap.String("replyMessage", replyMessage))
	return err
}
