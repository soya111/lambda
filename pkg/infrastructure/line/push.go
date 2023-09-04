package line

import (
	"context"
	"notify/pkg/logging"

	"github.com/hashicorp/go-multierror"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

// PushMessages push messages to LINE
func (b *Linebot) PushMessages(ctx context.Context, to []string, messages ...linebot.SendingMessage) error {
	var result *multierror.Error
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start pushing messages to LINE.", zap.Any("messages", messages))
	for _, recipient := range to {
		_, err := b.Client.PushMessage(recipient, messages...).WithContext(ctx).Do()
		if err != nil {
			result = multierror.Append(result, err)
			logger.Error("Failed to push messages to LINE.", zap.Error(err), zap.String("recipient", recipient))
		} else {
			logger.Info("Successfully pushed messages to LINE.", zap.String("recipient", recipient))
		}
	}
	if result != nil {
		logger.Warn("Failed to push messages to LINE for some users.")
	} else {
		logger.Info("All messages successfully pushed.")
	}
	return result.ErrorOrNil()
}
