// command_handler.go
package webhook

import (
	"context"
	"fmt"
	"notify/pkg/blog"
	"notify/pkg/infrastructure/line"
	"notify/pkg/model"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Command string
type CommandInfo struct {
	Handler CommandHandler
	Desc    string
}

type CommandHandler func(*linebot.Event, []string) error
type CommandHandlers map[Command]CommandInfo

const (
	RegCommand    Command = "reg"
	UnregCommand  Command = "unreg"
	ListCommand   Command = "list"
	WhoamiCommand Command = "whoami"
	HelpCommand   Command = "help"
	BlogCommand   Command = "blog"
	// 新しいコマンドを追加する場合はここに定義する
)

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
		BlogCommand: {
			Handler: h.handleBlogCommand,
			Desc:    "Get the latest blog of a member. Usage: blog [member]",
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
	if !model.IsMember(params[1]) {
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
	if !model.IsMember(params[1]) {
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
	var replyTextBuilder strings.Builder
	for command, cmdInfo := range h.getCommandHandlers() {
		replyTextBuilder.WriteString(fmt.Sprintf("%s: %s\n", string(command), cmdInfo.Desc))
	}

	// 最後の改行を取り除く
	replyText := replyTextBuilder.String()
	replyText = strings.TrimSuffix(replyText, "\n")

	if err := h.bot.ReplyMessage(context.TODO(), event.ReplyToken, linebot.NewTextMessage(replyText)); err != nil {
		return err
	}
	return nil
}

func (h *Handler) handleBlogCommand(event *linebot.Event, params []string) error {
	if len(params) < 2 {
		return fmt.Errorf("Member name must be provided. Usage: blog [member]")
	}

	memberName := params[1]
	if !model.IsMember(memberName) {
		if err := h.bot.ReplyTextMessages(context.TODO(), event.ReplyToken, fmt.Sprintf("%sは存在しません。", memberName)); err != nil {
			return fmt.Errorf("handleBlogCommand: %w", err)
		}
	}

	scraper := blog.NewHinatazakaScraper(nil)
	diary, err := scraper.GetLatestDiaryByMember(memberName)
	if err != nil {
		return h.bot.ReplyWithError(context.TODO(), event.ReplyToken, "内部エラー", err)
	}

	message := h.bot.CreateFlexMessage(diary)

	err = h.bot.ReplyMessage(context.TODO(), event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("handleBlogCommand: %w", err)
	}
	return nil
}

// type Subscriber struct {
// 	MemberName string `dynamo:"member_name" json:"member_name"  index:"user_id-index,range"`
// 	UserId     string `json:"user_id" dynamo:"user_id" index:"user_id-index,hash"`
// }

func (h *Handler) registerMember(member string, event *linebot.Event) error {
	ctx := context.TODO()
	token := event.ReplyToken

	id := line.ExtractEventSourceIdentifier(event)
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

	id := line.ExtractEventSourceIdentifier(event)
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

	id := line.ExtractEventSourceIdentifier(event)
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
	return h.bot.ReplyTextMessages(context.TODO(), event.ReplyToken, line.ExtractEventSourceIdentifier(event))
}
