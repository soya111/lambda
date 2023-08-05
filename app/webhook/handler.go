package webhook

import (
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

// type User struct {
// 	Id   string `json:"user_id" dynamodbav:"user_id"`
// 	Name string `json:"name" dynamodbav:"name"`
// }

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
