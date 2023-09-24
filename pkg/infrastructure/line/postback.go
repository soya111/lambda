package line

import (
	"encoding/json"
	"fmt"
	"zephyr/pkg/model"

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
	PostbackActionBlog       PostbackAction = "blog"
	PostbackActionProfile    PostbackAction = "prof"
	PostbackActionNickname   PostbackAction = "name"
	PostbackActionSelect     PostbackAction = "select"
)

const MemberKey = "member"
const ActionKey = "action"

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
	BlogLabel        = "最新のブログ"
	ProfileLabel     = "プロフィール"
	NicknameLabel    = "ニックネーム"
)

// NewSubscribeAction is the postback action that registers a member.
func NewSubscribeAction(diaryMemberName string) *linebot.PostbackAction {
	postBackMap := map[string]string{
		MemberKey: model.NormalizeName(diaryMemberName),
	}
	dataString, err := NewPostbackDataString(PostbackActionRegister, postBackMap)
	if err != nil {
		fmt.Printf("NewSubscribeAction: %v\n", err)
		return nil
	}
	return NewPostbackAction(SubscribeLabel, dataString, SubscribeLabel)
}

// newUnsubscribeAction is the postback action that unregisters a member.
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

// NewBlogAction is the postback action that shows the latest blog entry of the specified member.
func NewBlogAction(diaryMemberName string) *linebot.PostbackAction {
	postBackMap := map[string]string{
		MemberKey: model.NormalizeName(diaryMemberName),
	}
	dataString, err := NewPostbackDataString(PostbackActionBlog, postBackMap)
	if err != nil {
		fmt.Printf("NewBlogAction: %v\n", err)
		return nil
	}
	return NewPostbackAction(BlogLabel, dataString, BlogLabel)
}

// NewProfileAction is the postback action that shows the profile of the specified member.
func NewProfileAction(diaryMemberName string) *linebot.PostbackAction {
	postBackMap := map[string]string{
		MemberKey: model.NormalizeName(diaryMemberName),
	}
	dataString, err := NewPostbackDataString(PostbackActionProfile, postBackMap)
	if err != nil {
		fmt.Printf("NewProfileAction: %v\n", err)
		return nil
	}
	return NewPostbackAction(ProfileLabel, dataString, ProfileLabel)
}

// NewNicknameAction is the postback action that shows the nickname of the specified member.
func NewNicknameAction(diaryMemberName string) *linebot.PostbackAction {
	postBackMap := map[string]string{
		MemberKey: model.NormalizeName(diaryMemberName),
	}
	dataString, err := NewPostbackDataString(PostbackActionNickname, postBackMap)
	if err != nil {
		fmt.Printf("NewNicknameAction: %v\n", err)
		return nil
	}
	return NewPostbackAction(NicknameLabel, dataString, NicknameLabel)
}

// NewSelectAction is the postback action that shows the selectmenu of the member.
func NewSelectAction(label string) *linebot.PostbackAction {
	postBackMap := map[string]string{
		ActionKey: label,
	}
	dataString, err := NewPostbackDataString(PostbackActionSelect, postBackMap)
	if err != nil {
		fmt.Printf("NewSelectAction: %v\n", err)
		return nil
	}
	return NewPostbackAction(label, dataString, label)
}
