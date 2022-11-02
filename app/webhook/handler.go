package webhook

import (
	"context"
	"fmt"
	"strings"

	"notify/pkg/line"

	"github.com/guregu/dynamo"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Handler struct {
	bot *linebot.Client
	db  *dynamo.DB
}

func NewHandler(client *linebot.Client, db *dynamo.DB) *Handler {
	return &Handler{client, db}
}

func (h *Handler) HandleEvent(ctx context.Context, event *linebot.Event) error {
	switch event.Type {
	case linebot.EventTypeMessage:
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			err := h.handleTextMessage(message.Text, event)
			if err != nil {
				return fmt.Errorf("HandleEvent: %w", err)
			}
			return nil
		}
	case linebot.EventTypeLeave:
		// 退会、ブロック等されたら登録情報削除すべき
		// webhook.HandleEventLeave(event)
	}
	return nil
}

func (h *Handler) handleTextMessage(t string, event *linebot.Event) error {
	text := strings.Split(t, " ")
	switch {
	case text[0] == "reg" && isMember(text[1]):
		member := text[1]
		err := h.registerMember(member, event)
		if err != nil {
			return fmt.Errorf("handleTextMessage: %w", err)
		}
		return nil

	case text[0] == "unreg" && isMember(text[1]):
		member := text[1]
		err := h.unregisterMember(member, event)
		if err != nil {
			return fmt.Errorf("handleTextMessage: %w", err)
		}
		return nil

	case text[0] == "reg" && text[1] == "list":
		err := h.showSubscribeList(event)
		if err != nil {
			return fmt.Errorf("handleTextMessage: %w", err)
		}
		return nil

	case text[0] == "whoami":
		switch {
		case event.Source.Type == linebot.EventSourceTypeUser:
			h.sendUserId(event)
		case event.Source.Type == linebot.EventSourceTypeGroup:
			h.sendGroupId(event)
		}
	case isMember(text[0]):
		// TODO: いつか機能追加
		// 最新のブログおくるとか
	}
	return nil
}

func isMember(text string) bool {
	for _, v := range memberList {
		if text == v {
			return true
		}
	}
	return false
}

type Subscriber struct {
	MemberName string `dynamo:"member_name" json:"member_name"  index:"user_id-index,range"`
	UserId     string `json:"user_id" dynamo:"user_id" index:"user_id-index,hash"`
}

func (h *Handler) registerMember(member string, event *linebot.Event) error {
	var id string
	token := event.ReplyToken

	if event.Source.Type == linebot.EventSourceTypeUser {
		// user名調査
		userId := event.Source.UserID
		userProfile, _ := h.bot.GetProfile(userId).Do()
		err := h.postUser(User{userId, userProfile.DisplayName})
		if err != nil {
			fmt.Println(err)
		}
		id = userId
	} else if event.Source.Type == linebot.EventSourceTypeGroup {
		id = event.Source.GroupID
	}

	err := h.postSubscriber(Subscriber{member, id})
	if err != nil {
		message := "登録できませんでした！"
		if _, err := h.bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
			return fmt.Errorf("registerMember: %w", err)
		}
		return fmt.Errorf("registerMember: %w", err)
	}

	message := fmt.Sprintf("registered %s", member)
	if _, err := h.bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
		return fmt.Errorf("registerMember: %w", err)
	}
	return nil
}

func (h *Handler) unregisterMember(member string, event *linebot.Event) error {
	token := event.ReplyToken
	id := extractEventSourceIdentifier(event)

	err := h.deleteSubscriber(member, id)
	if err != nil {
		message := "登録できませんでした！"
		if _, err := h.bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
			return fmt.Errorf("unregisterMember: %w", err)
		}
		return fmt.Errorf("unregisterMember: %w", err)
	}

	message := fmt.Sprintf("unregistered %s", member)
	if _, err := h.bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
		return fmt.Errorf("unregisterMember: %w", err)
	}
	return nil
}

func (h *Handler) showSubscribeList(event *linebot.Event) error {
	token := event.ReplyToken
	id := extractEventSourceIdentifier(event)

	list, err := h.getSubscribeList(id)
	if err != nil {
		message := "情報を取得できませんでした！"
		if _, err := h.bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
			return fmt.Errorf("showSubscribeList: %w", err)
		}
		return fmt.Errorf("showSubscribeList: %w", err)

	}

	message := "登録リスト"
	for _, v := range list {
		message += fmt.Sprintf("\n%s", v.MemberName)
	}
	if _, err := h.bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
		return fmt.Errorf("showSubscribeList: %w", err)
	}
	return nil
}

func (h *Handler) postSubscriber(subscriber Subscriber) error {
	table := h.db.Table("Subscriber")

	if err := table.Put(subscriber).Run(); err != nil {
		return fmt.Errorf("postSubscriber: %w", err)
	}

	return nil
}

func (h *Handler) deleteSubscriber(memberName, userId string) error {
	table := h.db.Table("Subscriber")

	err := table.Delete("member_name", memberName).Range("user_id", userId).Run()

	if err != nil {
		return fmt.Errorf("deleteSubscriber: %w", err)
	}
	return nil
}

func (h *Handler) getSubscribeList(id string) ([]Subscriber, error) {
	table := h.db.Table("Subscriber")

	var res []Subscriber
	err := table.Get("user_id", id).Index("user_id-index").All(&res)
	if err != nil {
		return nil, fmt.Errorf("getSubscribeList: %w", err)
	}
	return res, nil
}

func (h *Handler) sendUserId(event *linebot.Event) {
	message := fmt.Sprintf("User id is \"%s\"", event.Source.UserID)
	line.NewLinebot().ReplyTextMessages(event.ReplyToken, message)
}

func (h *Handler) sendGroupId(event *linebot.Event) {
	message := fmt.Sprintf("Group id is \"%s\"\nYour user id is \"%s\"", event.Source.GroupID, event.Source.UserID)
	line.NewLinebot().ReplyTextMessages(event.ReplyToken, message)
}

type User struct {
	Id   string `json:"user_id" dynamodbav:"user_id"`
	Name string `json:"name" dynamodbav:"name"`
}

func (h *Handler) postUser(user User) error {
	table := h.db.Table("User")

	err := table.Put(user).Run()
	if err != nil {
		return fmt.Errorf("postUser: %w", err)
	}

	return nil
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
