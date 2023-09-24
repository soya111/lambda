package webhook

import (
	"context"
	"fmt"
	"zephyr/pkg/infrastructure/line"
	"zephyr/pkg/logging"
	"zephyr/pkg/model"
	"zephyr/pkg/service"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type PostbackCommand interface {
	Execute(context.Context, *linebot.Event, *line.PostbackData) error
}

type PostbackCommandMap map[line.PostbackAction]PostbackCommand

func (h *Handler) getPostbackCommandMap() PostbackCommandMap {
	subscriptionService := service.NewSubscriptionService(h.bot, h.subscriber)
	return PostbackCommandMap{
		line.PostbackActionRegister:   &PostbackCommandRegister{subscriptionService},
		line.PostbackActionUnregister: &PostbackCommandUnregister{subscriptionService},
	}
}

func (h *Handler) handlePostbackEvent(ctx context.Context, event *linebot.Event) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start handling postback event")

	data, err := line.ParsePostbackData(event)
	if err != nil {
		return fmt.Errorf("handlePostbackEvent: %w", err)
	}

	command, ok := h.getPostbackCommandMap()[data.Action]
	if !ok {
		return fmt.Errorf("unknown postback action: %s", data.Action)
	}
	return command.Execute(ctx, event, data)
}

// PostbackCommandRegister is a command to register a member.
type PostbackCommandRegister struct {
	subscriptionService *service.SubscriptionService
}

func (c *PostbackCommandRegister) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command register")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	err := c.subscriptionService.RegisterMember(ctx, member, event)
	if err != nil {
		return fmt.Errorf("PostbackCommandRegister.Execute: %w", err)
	}
	return nil
}

// PostbackCommandUnregister is a command to unregister a member.
type PostbackCommandUnregister struct {
	subscriptionService *service.SubscriptionService
}

func (c *PostbackCommandUnregister) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command unregister")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	err := c.subscriptionService.UnregisterMember(ctx, member, event)
	if err != nil {
		return fmt.Errorf("PostbackCommandUnregister.Execute: %w", err)
	}
	return nil
}
