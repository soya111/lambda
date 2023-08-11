package webhook

import (
	"notify/pkg/infrastructure/line"
	"notify/pkg/model"
)

// Handler is a struct that handles webhook events.
type Handler struct {
	bot        *line.Linebot
	subscriber model.SubscriberRepository
}

// NewHandler creates a new Handler.
func NewHandler(client *line.Linebot, subscriber model.SubscriberRepository) *Handler {
	return &Handler{client, subscriber}
}

// type User struct {
// 	Id   string `json:"user_id" dynamodbav:"user_id"`
// 	Name string `json:"name" dynamodbav:"name"`
// }
