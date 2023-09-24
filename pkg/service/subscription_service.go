package service

import (
	"context"
	"fmt"
	"zephyr/pkg/infrastructure/line"
	"zephyr/pkg/logging"
	"zephyr/pkg/model"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

// SubscriptionService is the struct that represents the subscription service.
type SubscriptionService struct {
	bot        *line.Linebot
	subscriber model.SubscriberRepository
}

// NewSubscriptionService creates a new SubscriptionService.
func NewSubscriptionService(bot *line.Linebot, subscriber model.SubscriberRepository) *SubscriptionService {
	return &SubscriptionService{bot, subscriber}
}

// RegisterMember registers the member to the subscription list.
func (s *SubscriptionService) RegisterMember(ctx context.Context, member string, event *linebot.Event) error {
	token := event.ReplyToken

	id := line.ExtractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return s.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	err := s.subscriber.Subscribe(model.Subscriber{MemberName: member, UserId: id})
	if err != nil {
		return s.bot.ReplyWithError(ctx, token, "登録できませんでした！", err)
	}

	logger := logging.LoggerFromContext(ctx)
	logger.Info("Registered member", zap.String("member", member), zap.String("id", id))

	message := fmt.Sprintf("registered %s", member)
	if err := s.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("registerMember: partial success, registration succeeded but failed to send message: %w", err)
	}
	return nil
}

// UnregisterMember unregisters the member from the subscription list.
func (s *SubscriptionService) UnregisterMember(ctx context.Context, member string, event *linebot.Event) error {
	token := event.ReplyToken

	id := line.ExtractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return s.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	err := s.subscriber.Unsubscribe(member, id)
	if err != nil {
		return s.bot.ReplyWithError(ctx, token, "登録解除できませんでした！", err)
	}

	logger := logging.LoggerFromContext(ctx)
	logger.Info("Unregistered member", zap.String("member", member), zap.String("id", id))

	message := fmt.Sprintf("unregistered %s", member)
	if err := s.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("unregisterMember: partial success, unregistration succeeded but failed to send message: %w", err)
	}
	return nil
}

// ShowSubscribeList shows the subscription list.
func (s *SubscriptionService) ShowSubscribeList(ctx context.Context, event *linebot.Event) error {
	token := event.ReplyToken

	id := line.ExtractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return s.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	list, err := s.subscriber.GetAllById(id)
	if err != nil {
		return s.bot.ReplyWithError(ctx, token, "情報を取得できませんでした！", err)
	}

	message := "登録リスト"
	for _, v := range list {
		message += fmt.Sprintf("\n%s", v.MemberName)
	}
	if err := s.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("showSubscribeList: %w", err)
	}

	logger := logging.LoggerFromContext(ctx)
	logger.Info("Showed subscribe list", zap.String("id", id))

	return nil
}
