package webhook

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"zephyr/pkg/blog"
	"zephyr/pkg/infrastructure/line"
	"zephyr/pkg/logging"
	"zephyr/pkg/model"
	"zephyr/pkg/profile"
	"zephyr/pkg/service"

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

const (
	CmdReg       CommandName = "reg"
	CmdUnreg     CommandName = "unreg"
	CmdList      CommandName = "list"
	CmdWhoami    CommandName = "whoami"
	CmdHelp      CommandName = "help"
	CmdBlog      CommandName = "blog"
	CmdProf      CommandName = "prof"
	CmdNickaname CommandName = "name"
	CmdMenu      CommandName = "menu"
	CmdId        CommandName = "id"
	// 新しいコマンドを追加する場合はここに定義する
)

func (h *Handler) getCommandHandlers() CommandMap {
	subscriptionService := service.NewSubscriptionService(h.bot, h.subscriber)
	identityService := service.NewIdentityService(h.bot)
	cmdMap := CommandMap{
		CmdReg:       &RegCommand{subscriptionService},
		CmdUnreg:     &UnregCommand{subscriptionService},
		CmdList:      &ListCommand{subscriptionService},
		CmdWhoami:    &WhoamiCommand{identityService},
		CmdBlog:      &BlogCommand{h.bot},
		CmdProf:      &ProfCommand{h.bot},
		CmdNickaname: &NicknameCommand{h.bot},
		CmdMenu:      &MenuCommand{h.bot},
		CmdId:        &IdCommand{h.bot},
		// 新たに追加するコマンドも同様にここに追加します
	}
	cmdMap[CmdHelp] = &HelpCommand{h.bot, cmdMap}
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
	subscriptionService *service.SubscriptionService
}

func (c *RegCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing RegCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}
	member := model.TranslateNicknametoMember(args[1])
	if !model.IsMember(member) {
		return nil
	}
	err := c.subscriptionService.RegisterMember(ctx, member, event)
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
	subscriptionService *service.SubscriptionService
}

func (c *UnregCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing UnregCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}
	member := model.TranslateNicknametoMember(args[1])
	if !model.IsMember(member) {
		return nil
	}
	err := c.subscriptionService.UnregisterMember(ctx, member, event)
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
	subscriptionService *service.SubscriptionService
}

func (c *ListCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing ListCommand")

	err := c.subscriptionService.ShowSubscribeList(ctx, event)
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
	identityService *service.IdentityService
}

func (c *WhoamiCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing WhoamiCommand")

	return c.identityService.SendWhoami(ctx, event)
}

func (c *WhoamiCommand) Description() string {
	return "Show your user or group ID. Usage: whoami"
}

// HelpCommand is the command that shows the list of available commands.
type HelpCommand struct {
	bot      *line.Linebot
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
	bot *line.Linebot
}

func (c *BlogCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing BlogCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}

	member := model.TranslateNicknametoMember(args[1])
	if model.IsGrad(member) {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sは卒業メンバーです。", member)); err != nil {
			return fmt.Errorf("BlogCommand.Execute: %w", err)
		}
		return nil
	}

	if !model.IsMember(member) {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sは存在しません。", member)); err != nil {
			return fmt.Errorf("BlogCommand.Execute: %w", err)
		}
		return nil
	}

	scraper := blog.NewHinatazakaScraper()
	diary, err := scraper.GetLatestDiaryByMember(member)
	if err != nil {
		return c.bot.ReplyWithError(ctx, event.ReplyToken, "内部エラー", err)
	}

	message := line.CreateFlexMessage(diary)

	err = c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("BlogCommand.Execute: %w", err)
	}
	return nil
}

func (c *BlogCommand) Description() string {
	return "Get the latest blog of a member. Usage: blog [member]"
}

// ProfCommand is the command that shows the profile of the specified member.
type ProfCommand struct {
	bot *line.Linebot
}

func (c *ProfCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing ProfCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}

	member := model.TranslateNicknametoMember(args[1])
	if model.IsGrad(member) {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sは卒業メンバーです。", member)); err != nil {
			return fmt.Errorf("ProfCommand.Execute: %w", err)
		}
		return nil
	}

	if !model.IsMember(member) {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sは存在しません。", member)); err != nil {
			return fmt.Errorf("ProfCommand.Execute: %w", err)
		}
		return nil
	}

	prof, _ := profile.ScrapeProfile(member)

	message := line.CreateProfileFlexMessage(prof)

	err := c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("ProfCommand.Execute: %w", err)
	}
	return nil
}

func (c *ProfCommand) Description() string {
	return "Get the profile of a member. Usage: prof [member]"
}

// NicknameCommand is the command that shows the nickname of the specified member.
type NicknameCommand struct {
	bot *line.Linebot
}

func (c *NicknameCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing NicknameCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}

	member := model.TranslateNicknametoMember(args[1])
	if model.IsGrad(member) {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sは卒業メンバーです。", member)); err != nil {
			return fmt.Errorf("NicknameCommand.Execute: %w", err)
		}
		return nil
	}

	if !model.IsMember(member) {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sは存在しません。", member)); err != nil {
			return fmt.Errorf("NicknameCommand.Execute: %w", err)
		}
		return nil
	}

	if member == model.Poka {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sにニックネームはありません。", member)); err != nil {
			return fmt.Errorf("NicknameCommand.Execute: %w", err)
		}
		return nil
	}

	prof, _ := profile.ScrapeProfile(member)

	message := line.CreateNicknameListFlexMessage(prof)

	err := c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("NicknameCommand.Execute: %w", err)
	}
	return nil
}

func (c *NicknameCommand) Description() string {
	return "Get the nickname of a member. Usage: name [member]"
}

// MenuCommand is the command that shows the menu of member selection, or the menu of the specified member if accompanied by the member name.
type MenuCommand struct {
	bot *line.Linebot
}

func (c *MenuCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing MenuCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		message := line.CreateMenuFlexMessage()

		err := c.bot.ReplyMessage(ctx, event.ReplyToken, message)
		if err != nil {
			return fmt.Errorf("NicknameCommand.Execute: %w", err)
		}
		return nil
	}

	member := model.TranslateNicknametoMember(args[1])
	if model.IsGrad(member) {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sは卒業メンバーです。", member)); err != nil {
			return fmt.Errorf("NicknameCommand.Execute: %w", err)
		}
		return nil
	}

	if !model.IsMember(member) {
		return nil
	}

	prof, _ := profile.ScrapeProfile(member)
	message := line.CreateMemberMenuFlexMessage(prof)

	err := c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("NicknameCommand.Execute: %w", err)
	}
	return nil
}

func (c *MenuCommand) Description() string {
	return "Get the general menu or the specified member menu. Usage: menu or menu [member]"
}

// IdCommandis the command that shows the specific blog from blog id.
type IdCommand struct {
	bot *line.Linebot
}

func (c *IdCommand) Execute(ctx context.Context, event *linebot.Event, args []string) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Executing IdCommand with args", zap.Any("args", args))

	if len(args) < 2 {
		return nil
	}

	blogId := model.TranslateNicknametoMember(args[1])

	pattern := `^\d{5}$`
	regex := regexp.MustCompile(pattern)

	if !regex.MatchString(blogId) {
		return nil
	}

	scraper := blog.NewHinatazakaScraper()
	diary, err := scraper.GetSpecificDiaryById(blogId)

	if err != nil {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, "無効なブログIDです。"); err != nil {
			return fmt.Errorf("NicknameCommand.Execute: %w", err)
		}
		return nil
	}

	message := line.CreateFlexMessage(diary)

	err = c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("IdCommand.Execute: %w", err)
	}
	return nil
}

func (c *IdCommand) Description() string {
	return "Get the specific blog from blog id. Usage: id [blogId]"
}

// type Subscriber struct {
// 	MemberName string `dynamo:"member_name" json:"member_name"  index:"user_id-index,range"`
// 	UserId     string `json:"user_id" dynamo:"user_id" index:"user_id-index,hash"`
// }
