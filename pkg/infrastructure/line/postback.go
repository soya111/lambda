package line

import (
	"encoding/json"
	"fmt"
	"notify/pkg/model"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// PostbackData is the struct that represents the postback data.
type PostbackData struct {
	Action PostbackAction    `json:"a"`
	Params map[string]string `json:"p,omitempty"`
}

type PostbackAction string

const (
	PostbackActionRegister   PostbackAction = "reg"
	PostbackActionUnregister PostbackAction = "unreg"
)

const MemberKey = "member"

// ParsePostbackData parses the postback data.
func ParsePostbackData(event *linebot.Event) (*PostbackData, error) {
	var data PostbackData
	err := json.Unmarshal([]byte(event.Postback.Data), &data)
	if err != nil {
		return nil, fmt.Errorf("ParsePostbackData: %w", err)
	}
	return &data, nil
}

// NewPostbackDataString returns the string of postback data.
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

// NewPostbackAction returns the postback action.
func NewPostbackAction(label, data, displayText string) *linebot.PostbackAction {
	return linebot.NewPostbackAction(label, data, "", displayText, "", "")
}

const (
	ThumbUpLabel     = "👍"
	ThumbDownLabel   = "👎"
	SubscribeLabel   = "購読する"
	UnsubscribeLabel = "解除する"
)

func NewSubscribeAction(diaryMemberName string) *linebot.PostbackAction {
	postBackMap := map[string]string{
		MemberKey: model.NormalizeName(diaryMemberName),
	}
	dataString, err := NewPostbackDataString(PostbackActionRegister, postBackMap)
	if err != nil {
		fmt.Printf("newSubscribeAction: %v\n", err)
		return nil
	}
	return NewPostbackAction(SubscribeLabel, dataString, SubscribeLabel)
}

func newUnsubscribeAction(diaryMemberName string) *linebot.PostbackAction {
	postBackMap := map[string]string{
		MemberKey: model.NormalizeName(diaryMemberName),
	}
	dataString, err := NewPostbackDataString(PostbackActionUnregister, postBackMap)
	if err != nil {
		fmt.Printf("newUnsubscribeAction: %v\n", err)
		return nil
	}
	return NewPostbackAction(UnsubscribeLabel, dataString, UnsubscribeLabel)
}
