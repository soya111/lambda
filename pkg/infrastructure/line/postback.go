package line

import (
	"encoding/json"
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type PostbackData struct {
	Action PostbackAction    `json:"a"`
	Params map[string]string `json:"p,omitempty"`
}

type PostbackAction string

const (
	Postback               PostbackAction = "postback"
	PostbackActionRegister PostbackAction = "reg"
)

const MemberKey = "member"

func ParsePostbackData(event *linebot.Event) (*PostbackData, error) {
	var data PostbackData
	err := json.Unmarshal([]byte(event.Postback.Data), &data)
	if err != nil {
		return nil, fmt.Errorf("ParsePostbackData: %w", err)
	}
	return &data, nil
}

func NewPostbackDataString(action PostbackAction, params map[string]string) (string, error) {
	data := PostbackData{
		Action: action,
		Params: params,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(dataBytes), nil
}

func NewPostbackAction(label, data, displayText string) *linebot.PostbackAction {
	return linebot.NewPostbackAction(label, data, "", displayText, "", "")
}
