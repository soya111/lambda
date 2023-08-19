package webhook

import (
	"context"
	"fmt"
	"notify/pkg/infrastructure/line"
	"notify/pkg/logging"
	"notify/pkg/model"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type PostbackCommand interface {
	Execute(context.Context, *linebot.Event, *line.PostbackData) error
}

type PostbackCommandMap map[line.PostbackAction]PostbackCommand

func (h *Handler) getPostbackCommandMap() PostbackCommandMap {
	base := NewBaseCommand(h.bot, h.subscriber)
	return PostbackCommandMap{
		line.PostbackActionRegister:   &PostbackCommandRegister{base},
		line.PostbackActionUnregister: &PostbackCommandUnregister{base},
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
	*BaseCommand
}

func (c *PostbackCommandRegister) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command register")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	err := c.registerMember(ctx, member, event)
	if err != nil {
		return fmt.Errorf("PostbackCommandRegister.Execute: %w", err)
	}
	return nil
}

// PostbackCommandUnregister is a command to unregister a member.
type PostbackCommandUnregister struct {
	*BaseCommand
}

func (c *PostbackCommandUnregister) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command unregister")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	err := c.unregisterMember(ctx, member, event)
	if err != nil {
		return fmt.Errorf("PostbackCommandUnregister.Execute: %w", err)
	}
	return nil
}
