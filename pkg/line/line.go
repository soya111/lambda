package line

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

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

// line送信
func (b *Linebot) PushTextMessages(ctx context.Context, to []string, messages ...string) error {
	var result *multierror.Error

	for _, message := range messages {
		for _, to := range to {
			if _, err := b.Client.PushMessage(to, linebot.NewTextMessage(message)).WithContext(ctx).Do(); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	return result.ErrorOrNil()
}

func (b *Linebot) PushFlexImagesMessage(ctx context.Context, to []string, urls []string) error {
	var result *multierror.Error

	contents := []*linebot.BubbleContainer{}
	for _, url := range urls {
		content := &linebot.BubbleContainer{
			Type: linebot.FlexContainerTypeBubble,
			Size: linebot.FlexBubbleSizeTypeMicro,
			Body: &linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeImage,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Contents: []linebot.FlexComponent{
					&linebot.ImageComponent{
						Type:   linebot.FlexComponentTypeImage,
						URL:    url,
						Margin: linebot.FlexComponentMarginTypeNone,
					},
				},
				Action: &linebot.URIAction{
					URI: url,
				},
			},
		}
		contents = append(contents, content)
	}

	container := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: contents,
	}

	for _, to := range to {
		if _, err := b.Client.PushMessage(to, linebot.NewFlexMessage("新着ブログがあります", container)).WithContext(ctx).Do(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

func (b *Linebot) ReplyTextMessages(ctx context.Context, token string, message string) error {
	if _, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(message)).WithContext(ctx).Do(); err != nil {
		return err
	}
	return nil
}
