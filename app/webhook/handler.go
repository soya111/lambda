package webhook

import (
	"context"
	"fmt"
	"strings"

	"notify/pkg/line"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Handler struct {
	bot *linebot.Client
}

func NewHandler(client *linebot.Client) *Handler {
	return &Handler{client}
}

func (h *Handler) HandleEvent(ctx context.Context, event *linebot.Event) error {
	switch event.Type {
	case linebot.EventTypeMessage:
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			h.handleTextMessage(message.Text, event)
		}
	case linebot.EventTypeLeave:
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
	MemberName string `dynamodbav:"member_name" json:"member_name"`
	UserId     string `json:"user_id" dynamodbav:"user_id"`
}

func (h *Handler) registerMember(member string, event *linebot.Event) error {
	var id string
	token := event.ReplyToken

	if event.Source.Type == linebot.EventSourceTypeUser {
		// user名調査
		userId := event.Source.UserID
		userProfile, _ := h.bot.GetProfile(userId).Do()
		err := postUser(&User{userId, userProfile.DisplayName})
		if err != nil {
			fmt.Println(err)
		}
		id = userId
	} else if event.Source.Type == linebot.EventSourceTypeGroup {
		id = event.Source.GroupID
	}

	err := h.postSubscriber(&Subscriber{member, id})
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

func (h *Handler) postSubscriber(subscriber *Subscriber) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	db := dynamodb.New(sess)

	// attribute value作成
	inputAV, err := dynamodbattribute.MarshalMap(subscriber)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Subscriber"),
		Item:      inputAV,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) deleteSubscriber(memberName, userId string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	db := dynamodb.New(sess)

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String("Subscriber"),
		Key: map[string]*dynamodb.AttributeValue{
			"member_name": {
				S: aws.String(memberName),
			},
			"user_id": {
				S: aws.String(userId),
			},
		},
	}

	_, err = db.DeleteItem(params)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) getSubscribeList(id string) ([]Subscriber, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)

	params := &dynamodb.QueryInput{
		TableName:              aws.String("Subscriber"),
		IndexName:              aws.String("user_id-index"),
		KeyConditionExpression: aws.String("#user_id = :user_id"),
		ExpressionAttributeNames: map[string]*string{
			"#user_id": aws.String("user_id"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user_id": {
				S: aws.String(id),
			},
		},
	}

	result, err := db.Query(params)
	if err != nil {
		return nil, err
	}

	memberList := []Subscriber{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &memberList)
	if err != nil {
		return nil, err
	}

	return memberList, nil
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

func postUser(user *User) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	db := dynamodb.New(sess)

	// attribute value作成
	inputAV, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("User"),
		Item:      inputAV,
	}

	_, err = db.PutItem(input)
	if err != nil {
		return err
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
