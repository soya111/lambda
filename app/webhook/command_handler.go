package webhook

import (
	"context"
	"fmt"
	"notify/pkg/blog"
	"notify/pkg/infrastructure/line"
	"notify/pkg/logging"
	"notify/pkg/model"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

// Command is the interface that wraps the basic Execute method.
type Command interface {
	Execute(context.Context, *linebot.Event, []string) error
	Description() string
}

// CommandName is the type that represents the command name.
type CommandName string

// CommandMap is the type that represents the map of command name and command.
type CommandMap map[CommandName]Command

// BaseCommand is the base struct for all commands.
type BaseCommand struct {
	bot        *line.Linebot
	subscriber model.SubscriberRepository
}

// NewBaseCommand creates a new BaseCommand.
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

func (h *Handler) handleTextMessage(ctx context.Context, param string, event *linebot.Event) error {
	params := strings.Split(param, " ")
	commandName := CommandName(params[0])
	command, ok := h.getCommandHandlers()[commandName]
	if !ok {
		return nil
	}
	return command.Execute(ctx, event, params)
}

// RegCommand is the command that registers a member.
type RegCommand struct {
	*BaseCommand
}

func (c *RegCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing RegCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}
	member := args[1]
	if !model.IsMember(member) {
		return nil
	}
	err := c.registerMember(ctx, member, event)
	if err != nil {
		return fmt.Errorf("RegCommand.Execute: %w", err)
	}
	return nil
}

func (c *RegCommand) Description() string {
	return "Register a member. Usage: reg [member]"
}

// UnregCommand is the command that unregisters a member.
type UnregCommand struct {
	*BaseCommand
}

func (c *UnregCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing UnregCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}
	member := args[1]
	if !model.IsMember(member) {
		return nil
	}
	err := c.unregisterMember(ctx, member, event)
	if err != nil {
		return fmt.Errorf("UnregCommand.Execute: %w", err)
	}
	return nil
}

func (c *UnregCommand) Description() string {
	return "Unregister a member. Usage: unreg [member]"
}

// ListCommand is the command that shows the list of registered members.
type ListCommand struct {
	*BaseCommand
}

func (c *ListCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing ListCommand")

	err := c.showSubscribeList(ctx, event)
	if err != nil {
		return fmt.Errorf("ListCommand.Execute: %w", err)
	}
	return nil
}

func (c *ListCommand) Description() string {
	return "Show the list of registered members. Usage: list"
}

// WhoamiCommand is the command that shows the user or group ID.
type WhoamiCommand struct {
	*BaseCommand
}

func (c *WhoamiCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing WhoamiCommand")

	return c.sendWhoami(ctx, event)
}

func (c *WhoamiCommand) Description() string {
	return "Show your user or group ID. Usage: whoami"
}

// HelpCommand is the command that shows the list of available commands.
type HelpCommand struct {
	*BaseCommand
	handlers CommandMap
}

func (c *HelpCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing HelpCommand")

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

// BlogCommand is the command that shows the latest blog entry of the specified member.
type BlogCommand struct {
	*BaseCommand
}

func (c *BlogCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing BlogCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}

	member := args[1]
	if !model.IsMember(member) {
		if err := c.bot.ReplyTextMessages(context.TODO(), event.ReplyToken, fmt.Sprintf("%sは存在しません。", member)); err != nil {
			return fmt.Errorf("BlogCommand.Execute: %w", err)
		}
	}

	scraper := blog.NewHinatazakaScraper()
	diary, err := scraper.GetLatestDiaryByMember(member)
	if err != nil {
		return c.bot.ReplyWithError(context.TODO(), event.ReplyToken, "内部エラー", err)
	}

	message := line.CreateFlexMessage(diary)

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

func (c *BaseCommand) registerMember(ctx context.Context, member string, event *linebot.Event) error {
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

	logger := logging.LoggerFromContext(ctx)
	logger.Info("Registered member", zap.String("member", member), zap.String("id", id))

	message := fmt.Sprintf("registered %s", member)
	if err := c.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("registerMember: partial success, registration succeeded but failed to send message: %w", err)
	}
	return nil
}

func (c *BaseCommand) unregisterMember(ctx context.Context, member string, event *linebot.Event) error {
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

	logger := logging.LoggerFromContext(ctx)
	logger.Info("Unregistered member", zap.String("member", member), zap.String("id", id))

	message := fmt.Sprintf("unregistered %s", member)
	if err := c.bot.ReplyTextMessages(ctx, token, message); err != nil {
		return fmt.Errorf("unregisterMember: partial success, unregistration succeeded but failed to send message: %w", err)
	}
	return nil
}

func (c *BaseCommand) showSubscribeList(ctx context.Context, event *linebot.Event) error {
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

	logger := logging.LoggerFromContext(ctx)
	logger.Info("Showed subscribe list", zap.String("id", id))

	return nil
}

func (c *BaseCommand) sendWhoami(ctx context.Context, event *linebot.Event) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Sending whoami")
	return c.bot.ReplyTextMessages(ctx, event.ReplyToken, line.ExtractEventSourceIdentifier(event))
}
