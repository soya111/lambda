package line

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/hashicorp/go-multierror"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	MessageBlogUpdate = "新着ブログがあります"
)

func (b *Linebot) CreateTextMessages(messages ...string) []linebot.SendingMessage {
	var sendingMessages []linebot.SendingMessage
	for _, message := range messages {
		sendingMessages = append(sendingMessages, linebot.NewTextMessage(message))
	}
	return sendingMessages
}

func (b *Linebot) CreateFlexImagesMessage(urls []string) linebot.SendingMessage {
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

	return linebot.NewFlexMessage(MessageBlogUpdate, container)
}

func (b *Linebot) PushMessages(ctx context.Context, to []string, messages []linebot.SendingMessage) error {
	var result *multierror.Error
	var requestId string

	lambdaCtx, ok := lambdacontext.FromContext(ctx)
	if ok {
		requestId = lambdaCtx.AwsRequestID
	}

	for _, to := range to {
		_, err := b.Client.PushMessage(to, messages...).WithContext(ctx).Do()
		result = multierror.Append(result, err)
		fmt.Printf("RequestId: %s, Push Messages to %s, error: %v\n", requestId, to, err)
	}
	return result.ErrorOrNil()
}

func (b *Linebot) ReplyTextMessages(ctx context.Context, token string, message string) error {
	if _, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(message)).WithContext(ctx).Do(); err != nil {
		return err
	}
	return nil
}
