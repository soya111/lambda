package webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func HandleTextMessage(t string, event *linebot.Event) {
	text := strings.Split(t, " ")
	switch {
	case text[0] == "reg":
		switch {
		case event.Source.Type == linebot.EventSourceTypeUser:
			registerUser(event)
		case event.Source.Type == linebot.EventSourceTypeGroup:
			registerGroup(event)
		}
	case text[0] == "unreg":
		switch {
		case event.Source.Type == linebot.EventSourceTypeUser:
			unregisterUser(event)
		case event.Source.Type == linebot.EventSourceTypeGroup:
			unregisterGroup(event)
		}
	case text[0] == "whoami":
		switch {
		case event.Source.Type == linebot.EventSourceTypeUser:
			sendUserId(event)
		case event.Source.Type == linebot.EventSourceTypeGroup:
			sendGroupId(event)
		}

	}
}

func registerUser(event *linebot.Event) {
	err := godotenv.Load(".env")

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	userId := event.Source.UserID
	userProfile, _ := getUserProfile(userId)
	err = postUser(&User{userId, userProfile.DisplayName})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("error")).Do(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		return
	}
	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("registered")).Do(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func registerGroup(event *linebot.Event) {
	err := godotenv.Load(".env")

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = postUser(&User{event.Source.GroupID, ""})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("error")).Do(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		return
	}
	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("registered")).Do(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func unregisterUser(event *linebot.Event) {
	err := godotenv.Load(".env")

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = deleteUser(event.Source.UserID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("error")).Do(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		return
	}
	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("unregistered")).Do(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func unregisterGroup(event *linebot.Event) {
	err := godotenv.Load(".env")

	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	err = deleteUser(event.Source.GroupID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("error")).Do(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		return
	}
	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("unregistered")).Do(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func sendUserId(event *linebot.Event) {
	message := fmt.Sprintf("User id is \"%s\"", event.Source.UserID)
	line.ReplyTextMessages(event.ReplyToken, message)
}

func sendGroupId(event *linebot.Event) {
	message := fmt.Sprintf("Group id is \"%s\"\nYour user id is \"%s\"", event.Source.GroupID, event.Source.UserID)
	line.ReplyTextMessages(event.ReplyToken, message)
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

func deleteUser(id string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	db := dynamodb.New(sess)

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String("User"),
		Key: map[string]*dynamodb.AttributeValue{
			"user_id": {
				S: aws.String(id),
			},
		},
	}

	_, err = db.DeleteItem(params)
	if err != nil {
		return err
	}

	return nil
}

type UserProfile struct {
	UserId        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureUrl    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	Language      string `json:"language"`
}

func getUserProfile(userId string) (*UserProfile, error) {
	err := godotenv.Load("../pkg/.env")

	url := fmt.Sprintf("https://api.line.me/v2/bot/profile/%s", userId)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CHANNEL_TOKEN")))
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	defer res.Body.Close()

	userProfile := &UserProfile{}
	err = json.NewDecoder(res.Body).Decode(userProfile)
	if err != nil {
		return nil, err
	}

	return userProfile, nil
}
