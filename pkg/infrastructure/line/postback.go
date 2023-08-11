package line

import (
	"encoding/json"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type PostbackData struct {
	Action PostbackAction    `json:"a"`
	Params map[string]string `json:"p,omitempty"`
}

type PostbackAction string

const (
	Postback                       PostbackAction = "postback"
	PostbackActionAddSubscriber    PostbackAction = "add_subscriber"
	PostbackActionRemoveSubscriber PostbackAction = "remove_subscriber"
	PostbackActionListSubscriber   PostbackAction = "list_subscriber"
	PostbackActionWhoami           PostbackAction = "whoami"
	PostbackActionHelp             PostbackAction = "help"
	PostbackActionBlog             PostbackAction = "blog"
)

func ParsePostbackData(event *linebot.Event) (*PostbackData, error) {
	var data PostbackData
	err := json.Unmarshal([]byte(event.Postback.Data), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
