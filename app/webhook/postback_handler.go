package webhook

import (
	"context"
	"fmt"
	"zephyr/pkg/blog"
	"zephyr/pkg/infrastructure/line"
	"zephyr/pkg/logging"
	"zephyr/pkg/model"
	"zephyr/pkg/profile"
	"zephyr/pkg/service"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type PostbackCommand interface {
	Execute(context.Context, *linebot.Event, *line.PostbackData) error
}

type PostbackCommandMap map[line.PostbackAction]PostbackCommand

func (h *Handler) getPostbackCommandMap() PostbackCommandMap {
	subscriptionService := service.NewSubscriptionService(h.bot, h.subscriber)
	return PostbackCommandMap{
		line.PostbackActionRegister:   &PostbackCommandRegister{subscriptionService},
		line.PostbackActionUnregister: &PostbackCommandUnregister{subscriptionService},
		line.PostbackActionBlog:       &PostbackCommandBlog{h.bot},
		line.PostbackActionProfile:    &PostbackCommandProfile{h.bot},
		line.PostbackActionNickname:   &PostbackCommandNickname{h.bot},
		line.PostbackActionSelect:     &PostbackCommandSelect{h.bot},
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
	subscriptionService *service.SubscriptionService
}

func (c *PostbackCommandRegister) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command register")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	err := c.subscriptionService.RegisterMember(ctx, member, event)
	if err != nil {
		return fmt.Errorf("PostbackCommandRegister.Execute: %w", err)
	}
	return nil
}

// PostbackCommandUnregister is a command to unregister a member.
type PostbackCommandUnregister struct {
	subscriptionService *service.SubscriptionService
}

func (c *PostbackCommandUnregister) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command unregister")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	err := c.subscriptionService.UnregisterMember(ctx, member, event)
	if err != nil {
		return fmt.Errorf("PostbackCommandUnregister.Execute: %w", err)
	}
	return nil
}

// PostbackCommandBlog is a command to show the latest blog entry of the specified member.
type PostbackCommandBlog struct {
	bot *line.Linebot
}

func (c *PostbackCommandBlog) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command blog")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	scraper := blog.NewHinatazakaScraper()
	diary, err := scraper.GetLatestDiaryByMember(member)
	if err != nil {
		return c.bot.ReplyWithError(ctx, event.ReplyToken, "内部エラー", err)
	}

	message := line.CreateFlexMessage(diary)

	err = c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("PostbackCommandBlog.Execute: %w", err)
	}
	return nil
}

// PostbackCommandProfile is a command to show the profile of the specified member.
type PostbackCommandProfile struct {
	bot *line.Linebot
}

func (c *PostbackCommandProfile) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command profile")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	prof, _ := profile.ScrapeProfile(member)

	message := line.CreateProfileFlexMessage(prof)

	err := c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("PostbackCommandProfile.Execute: %w", err)
	}
	return nil
}

// PostbackCommandNickname is a command to show the nickname of the specified member.
type PostbackCommandNickname struct {
	bot *line.Linebot
}

func (c *PostbackCommandNickname) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command nickname")

	member := data.Params[line.MemberKey]
	if !model.IsMember(member) {
		return fmt.Errorf("invalid member: %s", member)
	}

	if member == model.Poka {
		if err := c.bot.ReplyTextMessages(ctx, event.ReplyToken, fmt.Sprintf("%sにニックネームはありません。", member)); err != nil {
			return fmt.Errorf("PostbackCommandNickname.Execute: %w", err)
		}
		return nil
	}

	prof, _ := profile.ScrapeProfile(member)

	message := line.CreateNicknameListFlexMessage(prof)

	err := c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("PostbackCommandNickname.Execute: %w", err)
	}
	return nil
}

// PostbackCommandSelect is a command to show the selectmenu of the member.
type PostbackCommandSelect struct {
	bot *line.Linebot
}

var labelToActionMap = map[string]line.PostbackActionGenerator{
	line.SubscribeLabel: line.NewSubscribeAction,
	line.BlogLabel:      line.NewBlogAction,
	line.ProfileLabel:   line.NewProfileAction,
	line.NicknameLabel:  line.NewNicknameAction,
}

func (c *PostbackCommandSelect) Execute(ctx context.Context, event *linebot.Event, data *line.PostbackData) error {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Start executing postback command select")

	label := data.Params[line.ActionKey]

	message := line.CreateMemberSelectFlexMessage(labelToActionMap[label])

	err := c.bot.ReplyMessage(ctx, event.ReplyToken, message)
	if err != nil {
		return fmt.Errorf("PostbackCommandSelect.Execute: %w", err)
	}

	return nil
}
