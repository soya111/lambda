package webhook

import (
	"context"
	"errors"
	"zephyr/pkg/infrastructure/line"
	"zephyr/pkg/logging"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type EventHandler func(ctx context.Context, event *linebot.Event) error

type EventHandlers map[linebot.EventType]EventHandler

func (h *Handler) getEventHandlers() EventHandlers {
	return EventHandlers{
		linebot.EventTypeMessage:  h.handleMessageEvent,
		linebot.EventTypeFollow:   h.handleFollowEvent,
		linebot.EventTypeUnfollow: h.handleUnfollowEvent,
		linebot.EventTypeJoin:     h.handleJoinEvent,
		linebot.EventTypeLeave:    h.handleLeaveEvent,
		linebot.EventTypePostback: h.handlePostbackEvent,
		// 他のイベントタイプも同様に定義する
	}
}

func (h *Handler) handleMessageEvent(ctx context.Context, event *linebot.Event) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Handling message event")

	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		return h.handleTextMessage(ctx, message.Text, event)
	default:
		// TextMessage以外は何もしない
		logger.Info("Unsupported message type")
		return nil
	}
}

func (h *Handler) handleFollowEvent(ctx context.Context, event *linebot.Event) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Handling follow event")
	return nil
}

func (h *Handler) handleUnfollowEvent(ctx context.Context, event *linebot.Event) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Handling unfollow event")

	id := line.ExtractEventSourceIdentifier(event)
	if id == "" {
		logger.Error("Failed to extract event source identifier")
		return errors.New("failed to extract event source identifier")
	}

	if err := h.subscriber.DeleteAllById(id); err != nil {
		logger.Error("Failed to delete subscriber", zap.Error(err))
		return err
	}

	logger.Info("Successfully handled unfollow event")

	return nil
}

func (h *Handler) handleJoinEvent(ctx context.Context, event *linebot.Event) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Handling join event")
	return nil
}

func (h *Handler) handleLeaveEvent(ctx context.Context, event *linebot.Event) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Handling leave event")

	id := line.ExtractEventSourceIdentifier(event)
	if id == "" {
		logger.Error("Failed to extract event source identifier")
		return errors.New("failed to extract event source identifier")
	}

	if err := h.subscriber.DeleteAllById(id); err != nil {
		logger.Error("Failed to delete subscriber", zap.Error(err))
		return err
	}

	logger.Info("Successfully handled leave event")

	return nil
}

// HandleEvent method remains in the main handler file or the file where it's being called.
func (h *Handler) HandleEvent(ctx context.Context, event *linebot.Event) error {
	handler, ok := h.getEventHandlers()[event.Type]
	if ok {
		return handler(ctx, event)
	}

	logger := logging.LoggerFromContext(ctx)
	logger.Error("Unknown event type", zap.String("EventType", string(event.Type)))

	return nil
}
