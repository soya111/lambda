package webhook

import (
	"context"
	"fmt"
	"notify/pkg/infrastructure/line"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type PostbackCommand interface {
	Execute(context.Context, *linebot.Event, *line.PostbackData) error
}

type PostbackCommandMap map[line.PostbackAction]PostbackCommand

func (h *Handler) getPostbackCommandMap() PostbackCommandMap {
	return PostbackCommandMap{
		line.Postback: nil,
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
