package webhook

import (
	"notify/pkg/infrastructure/line"
	"notify/pkg/model"
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
