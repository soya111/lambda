package webhook

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type EventHandler func(ctx context.Context, event *linebot.Event) error

type EventHandlers map[linebot.EventType]EventHandler

func (h *Handler) getEventHandlers() EventHandlers {
	return EventHandlers{
		linebot.EventTypeMessage:  h.handleMessageEvent,
		linebot.EventTypeLeave:    h.handleLeaveEvent,
		linebot.EventTypePostback: h.handlePostbackEvent,
		// 他のイベントタイプも同様に定義する
	}
}

func (h *Handler) handleMessageEvent(ctx context.Context, event *linebot.Event) error {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		return h.handleTextMessage(message.Text, event)
	default:
		// TextMessage以外は何もしない
		return nil
	}
}

func (h *Handler) handleLeaveEvent(ctx context.Context, event *linebot.Event) error {
	// EventTypeLeaveのときの処理を記述
	// ...
	return nil
}

// HandleEvent method remains in the main handler file or the file where it's being called.
func (h *Handler) HandleEvent(ctx context.Context, event *linebot.Event) error {
	handler, ok := h.getEventHandlers()[event.Type]
	if ok {
		return handler(ctx, event)
	}
	// 不明なイベントタイプに対するエラーハンドリング
	// ...
	return nil
}
