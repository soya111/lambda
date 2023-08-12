package webhook

import (
	"context"
	"fmt"
	"notify/pkg/infrastructure/line"
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
		line.PostbackActionRegister: &PostbackCommandRegister{base},
	}
}

func (h *Handler) handlePostbackEvent(ctx context.Context, event *linebot.Event) error {
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

type PostbackCommandRegister struct {
	*BaseCommand
}

func (c *PostbackCommandRegister) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	member := data.Params["member"]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	err := c.registerMember(member, event)
	if err != nil {
		return fmt.Errorf("PostbackCommandRegister.Execute: %w", err)
	}
	return nil
}
