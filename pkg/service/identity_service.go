package service

import (
	"context"
	"notify/pkg/infrastructure/line"
	"notify/pkg/logging"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// IdentityService is the struct that represents the identity service.
type IdentityService struct {
	bot *line.Linebot
}

// NewIdentityService creates a new IdentityService.
func NewIdentityService(bot *line.Linebot) *IdentityService {
	return &IdentityService{bot}
}

// SendWhoami sends the whoami message.
func (s *IdentityService) SendWhoami(ctx context.Context, event *linebot.Event) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Sending whoami")
	return s.bot.ReplyTextMessages(ctx, event.ReplyToken, line.ExtractEventSourceIdentifier(event))
}
