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

type Command interface {
	Execute(*linebot.Event, []string) error
	Description() string
}

type CommandName string

type CommandMap map[CommandName]Command

type BaseCommand struct {
	bot        *line.Linebot
	subscriber model.SubscriberRepository
}

func NewBaseCommand(bot *line.Linebot, subscriber model.SubscriberRepository) *BaseCommand {
	return &BaseCommand{bot, subscriber}
}

const (
	CmdReg    CommandName = "reg"
	CmdUnreg  CommandName = "unreg"
	CmdList   CommandName = "list"
	CmdWhoami CommandName = "whoami"
	CmdHelp   CommandName = "help"
	CmdBlog   CommandName = "blog"
	// 新しいコマンドを追加する場合はここに定義する
)

func (h *Handler) getCommandHandlers() CommandMap {
	base := NewBaseCommand(h.bot, h.subscriber)
	cmdMap := CommandMap{
		CmdReg:    &RegCommand{base},
		CmdUnreg:  &UnregCommand{base},
		CmdList:   &ListCommand{base},
		CmdWhoami: &WhoamiCommand{base},
		CmdBlog:   &BlogCommand{base},
		// 新たに追加するコマンドも同様にここに追加します
	}
	cmdMap[CmdHelp] = &HelpCommand{base, cmdMap}
	return cmdMap
}

func (h *Handler) handleTextMessage(param string, event *linebot.Event) error {
	params := strings.Split(param, " ")
	commandName := CommandName(params[0])
	command, ok := h.getCommandHandlers()[commandName]
	if !ok {
		return nil
	}
	return command.Execute(event, params)
}

type RegCommand struct {
	*BaseCommand
}

func (c *RegCommand) Execute(event *linebot.Event, args []string) error {
	if len(args) < 2 {
		return nil
	}
	member := args[1]
	if !model.IsMember(member) {
		return nil
	}
	err := c.registerMember(member, event)
	if err != nil {
		return fmt.Errorf("RegCommand.Execute: %w", err)
	}
	return nil
}

func (c *RegCommand) Description() string {
	return "Register a member. Usage: reg [member]"
}

type UnregCommand struct {
	*BaseCommand
}

func (c *UnregCommand) Execute(event *linebot.Event, args []string) error {
	if len(args) < 2 {
		return nil
	}
	member := args[1]
	if !model.IsMember(member) {
		return nil
	}
	err := c.unregisterMember(member, event)
	if err != nil {
		return fmt.Errorf("UnregCommand.Execute: %w", err)
	}
	return nil
}

func (c *UnregCommand) Description() string {
	return "Unregister a member. Usage: unreg [member]"
}

type ListCommand struct {
	*BaseCommand
}

func (c *ListCommand) Execute(event *linebot.Event, args []string) error {
	err := c.showSubscribeList(event)
	if err != nil {
		return fmt.Errorf("ListCommand.Execute: %w", err)
	}
	return nil
}

func (c *ListCommand) Description() string {
	return "Show the list of registered members. Usage: list"
}

type WhoamiCommand struct {
	*BaseCommand
}

func (c *WhoamiCommand) Execute(event *linebot.Event, args []string) error {
	return c.sendWhoami(event)
}

func (c *WhoamiCommand) Description() string {
	return "Show your user or group ID. Usage: whoami"
}

type HelpCommand struct {
	*BaseCommand
	handlers CommandMap
}

func (c *HelpCommand) Execute(event *linebot.Event, args []string) error {
	var replyTextBuilder strings.Builder
	cmdMap := c.handlers
	for commandName, command := range cmdMap {
		replyTextBuilder.WriteString(fmt.Sprintf("%s: %s\n", string(commandName), command.Description()))
	}

	// 最後の改行を取り除く
	replyText := replyTextBuilder.String()
	replyText = strings.TrimSuffix(replyText, "\n")

	if err := c.bot.ReplyMessage(context.TODO(), event.ReplyToken, linebot.NewTextMessage(replyText)); err != nil {
		return err
	}
	return nil
}

func (c *HelpCommand) Description() string {
	return "Show the list of available commands. Usage: help"
}

type BlogCommand struct {
	*BaseCommand
}

func (c *BlogCommand) Execute(event *linebot.Event, args []string) error {
	if len(args) < 2 {
		return nil
	}

	member := args[1]
	if !model.IsMember(member) {
		if err := c.bot.ReplyTextMessages(context.TODO(), event.ReplyToken, fmt.Sprintf("%sは存在しません。", member)); err != nil {
			return fmt.Errorf("BlogCommand.Execute: %w", err)
		}
	}

	scraper := blog.NewHinatazakaScraper(nil)
	diary, err := scraper.GetLatestDiaryByMember(member)
	if err != nil {
		return c.bot.ReplyWithError(context.TODO(), event.ReplyToken, "内部エラー", err)
	}

	message := c.bot.CreateFlexMessage(diary)

	err = c.bot.ReplyMessage(context.TODO(), event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("BlogCommand.Execute: %w", err)
	}
	return nil
}

func (c *BlogCommand) Description() string {
	return "Get the latest blog of a member. Usage: blog [member]"
}

// type Subscriber struct {
// 	MemberName string `dynamo:"member_name" json:"member_name"  index:"user_id-index,range"`
// 	UserId     string `json:"user_id" dynamo:"user_id" index:"user_id-index,hash"`
// }

func (c *RegCommand) registerMember(member string, event *linebot.Event) error {
	ctx := context.TODO()
	token := event.ReplyToken

	id := line.ExtractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return c.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	err := c.subscriber.Subscribe(model.Subscriber{MemberName: member, UserId: id})
	if err != nil {
		return c.bot.ReplyWithError(ctx, token, "登録できませんでした！", err)
	}

	message := fmt.Sprintf("registered %s", member)
	if err := c.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("registerMember: %w", err)
	}
	return nil
}

func (c *UnregCommand) unregisterMember(member string, event *linebot.Event) error {
	ctx := context.TODO() // あるいは他の適切なコンテキストを使用します
	token := event.ReplyToken

	id := line.ExtractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return c.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	err := c.subscriber.Unsubscribe(member, id)
	if err != nil {
		return c.bot.ReplyWithError(ctx, token, "登録解除できませんでした！", err)
	}

	message := fmt.Sprintf("unregistered %s", member)
	if err := c.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("unregisterMember: %w", err)
	}
	return nil
}

func (c *ListCommand) showSubscribeList(event *linebot.Event) error {
	ctx := context.TODO() // あるいは他の適切なコンテキストを使用します
	token := event.ReplyToken

	id := line.ExtractEventSourceIdentifier(event)
	if id == "" {
		err := fmt.Errorf("invalid source type: %v", event.Source.Type)
		return c.bot.ReplyWithError(ctx, token, "Invalid source type!", err)
	}

	list, err := c.subscriber.GetAllById(id)
	if err != nil {
		return c.bot.ReplyWithError(ctx, token, "情報を取得できませんでした！", err)
	}

	message := "登録リスト"
	for _, v := range list {
		message += fmt.Sprintf("\n%s", v.MemberName)
	}
	if err := c.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("showSubscribeList: %w", err)
	}
	return nil
}

func (c *WhoamiCommand) sendWhoami(event *linebot.Event) error {
	return c.bot.ReplyTextMessages(context.TODO(), event.ReplyToken, line.ExtractEventSourceIdentifier(event))
}
