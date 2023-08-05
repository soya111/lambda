package webhook

import (
	"context"
	"fmt"
	"strings"

	"notify/pkg/infrastructure/line"
	"notify/pkg/model"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Handler struct {
	bot        *line.Linebot
	subscriber model.SubscriberRepository
}

func NewHandler(client *line.Linebot, subscriber model.SubscriberRepository) *Handler {
	return &Handler{client, subscriber}
}

type EventHandler func(ctx context.Context, event *linebot.Event) error

type EventHandlers map[linebot.EventType]EventHandler

func (h *Handler) getEventHandlers() EventHandlers {
	return EventHandlers{
		linebot.EventTypeMessage: h.handleMessageEvent,
		linebot.EventTypeLeave:   h.handleLeaveEvent,
		// 他のイベントタイプも同様に定義する
	}
}

func (h *Handler) HandleEvent(ctx context.Context, event *linebot.Event) error {
	handler, ok := h.getEventHandlers()[event.Type]
	if ok {
		return handler(ctx, event)
	}
	// 不明なイベントタイプに対するエラーハンドリング
	// ...
	return nil
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

type Command string

const (
	RegCommand    Command = "reg"
	UnregCommand  Command = "unreg"
	ListCommand   Command = "list"
	WhoamiCommand Command = "whoami"
	HelpCommand   Command = "help"
	// 新しいコマンドを追加する場合はここに定義する
)

type CommandInfo struct {
	Handler CommandHandler
	Desc    string
}
type CommandHandler func(*linebot.Event, []string) error

type CommandHandlers map[Command]CommandInfo

func (h *Handler) getCommandHandlers() CommandHandlers {
	return CommandHandlers{
		RegCommand: {
			Handler: h.handleRegCommand,
			Desc:    "Register a member. Usage: reg [member]",
		},
		UnregCommand: {
			Handler: h.handleUnregCommand,
			Desc:    "Unregister a member. Usage: unreg [member]",
		},
		ListCommand: {
			Handler: h.handleListCommand,
			Desc:    "Show the list of registered members. Usage: list",
		},
		WhoamiCommand: {
			Handler: h.handleWhoamiCommand,
			Desc:    "Show your user or group ID. Usage: whoami",
		},
		HelpCommand: {
			Handler: h.handleHelpCommand,
			Desc:    "Show the list of available commands. Usage: help",
		},
		// 新たに追加するコマンドも同様にここに追加します
	}
}

func (h *Handler) handleTextMessage(param string, event *linebot.Event) error {
	params := strings.Split(param, " ")
	command := Command(params[0])
	cmdInfo, ok := h.getCommandHandlers()[command]
	if !ok {
		return nil
	}
	return cmdInfo.Handler(event, params)
}

func (h *Handler) handleRegCommand(event *linebot.Event, params []string) error {
	if len(params) < 2 {
		return nil
	}
	if !isMember(params[1]) {
		return nil
	}
	member := params[1]
	err := h.registerMember(member, event)
	if err != nil {
		return fmt.Errorf("handleRegCommand: %w", err)
	}
	return nil
}

func (h *Handler) handleUnregCommand(event *linebot.Event, params []string) error {
	if len(params) < 2 {
		return nil
	}
	if !isMember(params[1]) {
		return nil
	}
	member := params[1]
	err := h.unregisterMember(member, event)
	if err != nil {
		return fmt.Errorf("handleUnregCommand: %w", err)
	}
	return nil
}

func (h *Handler) handleListCommand(event *linebot.Event, params []string) error {
	err := h.showSubscribeList(event)
	if err != nil {
		return fmt.Errorf("handleListCommand: %w", err)
	}
	return nil
}

func (h *Handler) handleWhoamiCommand(event *linebot.Event, params []string) error {
	return h.sendWhoami(event)
}

func (h *Handler) handleHelpCommand(event *linebot.Event, params []string) error {
	replyText := "Available commands are:\n"
	for command, cmdInfo := range h.getCommandHandlers() {
		replyText += fmt.Sprintf("%s: %s\n", string(command), cmdInfo.Desc)
	}
	if _, err := h.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do(); err != nil {
		return err
	}
	return nil
}

type Subscriber struct {
	MemberName string `dynamo:"member_name" json:"member_name"  index:"user_id-index,range"`
	UserId     string `json:"user_id" dynamo:"user_id" index:"user_id-index,hash"`
}

func (h *Handler) registerMember(member string, event *linebot.Event) error {
	ctx := context.TODO()
	token := event.ReplyToken

	id := extractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return h.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	err := h.subscriber.Subscribe(model.Subscriber{MemberName: member, UserId: id})
	if err != nil {
		return h.bot.ReplyWithError(ctx, token, "登録できませんでした！", err)
	}

	message := fmt.Sprintf("registered %s", member)
	if err := h.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("registerMember: %w", err)
	}
	return nil
}

func (h *Handler) unregisterMember(member string, event *linebot.Event) error {
	ctx := context.TODO() // あるいは他の適切なコンテキストを使用します
	token := event.ReplyToken

	id := extractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return h.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	err := h.subscriber.Unsubscribe(member, id)
	if err != nil {
		return h.bot.ReplyWithError(ctx, token, "登録解除できませんでした！", err)
	}

	message := fmt.Sprintf("unregistered %s", member)
	if err := h.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("unregisterMember: %w", err)
	}
	return nil
}

func (h *Handler) showSubscribeList(event *linebot.Event) error {
	ctx := context.TODO() // あるいは他の適切なコンテキストを使用します
	token := event.ReplyToken

	id := extractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return h.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	list, err := h.subscriber.GetAllById(id)
	if err != nil {
		return h.bot.ReplyWithError(ctx, token, "情報を取得できませんでした！", err)
	}

	message := "登録リスト"
	for _, v := range list {
		message += fmt.Sprintf("\n%s", v.MemberName)
	}
	if err := h.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("showSubscribeList: %w", err)
	}
	return nil
}

func (h *Handler) sendWhoami(event *linebot.Event) error {
	return h.bot.ReplyTextMessages(context.TODO(), event.ReplyToken, extractEventSourceIdentifier(event))
}

type User struct {
	Id   string `json:"user_id" dynamodbav:"user_id"`
	Name string `json:"name" dynamodbav:"name"`
}

func extractEventSourceIdentifier(event *linebot.Event) string {
	var id string

	if event.Source.Type == linebot.EventSourceTypeUser {
		id = event.Source.UserID
	} else if event.Source.Type == linebot.EventSourceTypeGroup {
		id = event.Source.GroupID
	} else if event.Source.Type == linebot.EventSourceTypeRoom {
		id = event.Source.RoomID
	}

	return id
}
