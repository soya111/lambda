package line

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Client interface {
	PushTextMessages(ctx context.Context, to []string, messages ...string) error
	PushFlexImagesMessage(ctx context.Context, to []string, urls []string) error
	ReplyTextMessages(ctx context.Context, token string, message string) error
}

type Linebot struct {
	*linebot.Client
}

// 本番用コンストラクタ
func NewLinebot(channelSecret string, channelToken string) (*Linebot, error) {
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}
	return &Linebot{bot}, nil
}
