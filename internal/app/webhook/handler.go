package webhook

import (
	"context"
	"fmt"
	"os"
	"strings"

	"notify/internal/pkg/line"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Handler struct{}

func HandleEvent(ctx context.Context, event *linebot.Event) error {
	switch event.Type {
	case linebot.EventTypeMessage:
		switch message := event.Message.(type) {
		case *linebot.TextMessage:
			handleTextMessage(message.Text, event)
		}
	case linebot.EventTypeLeave:
		// webhook.HandleEventLeave(event)
	}
	return nil
}

func handleTextMessage(t string, event *linebot.Event) {
	text := strings.Split(t, " ")
	switch {
	case text[0] == "reg" && isMember(text[1]):
		member := text[1]
		registerMember(member, event)
	case text[0] == "unreg" && isMember(text[1]):
		member := text[1]
		unregisterMember(member, event)
	case text[0] == "reg" && text[1] == "list":
		showSubscribeList(event)
	case text[0] == "whoami":
		switch {
		case event.Source.Type == linebot.EventSourceTypeUser:
			sendUserId(event)
		case event.Source.Type == linebot.EventSourceTypeGroup:
			sendGroupId(event)
		}
	case isMember(text[0]):
		// TODO: いつか機能追加
	}
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

func registerMember(member string, event *linebot.Event) {
	err := godotenv.Load(".env")

	bot := line.NewLinebot()

	var id string

	if event.Source.Type == linebot.EventSourceTypeUser {
		// user名調査
		userId := event.Source.UserID
		userProfile, _ := bot.Client.GetProfile(userId).Do()
		_ = postUser(&User{userId, userProfile.DisplayName})
		id = userId
	} else if event.Source.Type == linebot.EventSourceTypeGroup {
		id = event.Source.GroupID
	}

	err = postSubscriber(&Subscriber{member, id})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err = bot.ReplyTextMessages(event.ReplyToken, "error"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		return
	}

	message := fmt.Sprintf("registered %s", member)
	if err = bot.ReplyTextMessages(event.ReplyToken, message); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func unregisterMember(member string, event *linebot.Event) {
	err := godotenv.Load(".env")

	bot := line.NewLinebot()

	id := extractEventSourceIdentifier(event)

	err = deleteSubscriber(member, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err = bot.ReplyTextMessages(event.ReplyToken, "error"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		return
	}

	message := fmt.Sprintf("unregistered %s", member)
	if err = bot.ReplyTextMessages(event.ReplyToken, message); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func showSubscribeList(event *linebot.Event) {
	err := godotenv.Load(".env")

	bot := line.NewLinebot()

	id := extractEventSourceIdentifier(event)

	list, err := getSubscribeList(id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err = bot.ReplyTextMessages(event.ReplyToken, "error"); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		return
	}

	message := "登録リスト"
	for _, v := range list {
		message += fmt.Sprintf("\n%s", v.MemberName)
	}
	if err = bot.ReplyTextMessages(event.ReplyToken, message); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func postSubscriber(subscriber *Subscriber) error {
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

func deleteSubscriber(memberName, userId string) error {
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

func getSubscribeList(id string) ([]Subscriber, error) {
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

func sendUserId(event *linebot.Event) {
	message := fmt.Sprintf("User id is \"%s\"", event.Source.UserID)
	line.NewLinebot().ReplyTextMessages(event.ReplyToken, message)
}

func sendGroupId(event *linebot.Event) {
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
