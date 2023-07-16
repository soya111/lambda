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

var MicroBubbleContainer = linebot.BubbleContainer{
	Type: linebot.FlexContainerTypeBubble,
	Size: linebot.FlexBubbleSizeTypeMicro,
}

func (b *Linebot) CreateTextMessages(messages ...string) []linebot.SendingMessage {
	var sendingMessages []linebot.SendingMessage
	for _, message := range messages {
		sendingMessages = append(sendingMessages, linebot.NewTextMessage(message))
	}
	return sendingMessages
}

func (b *Linebot) CreateFlexMessage(message string, urls []string) linebot.SendingMessage {
	var container []*linebot.BubbleContainer
	container = append(container, b.CreateFlexTextMessage(message))

	if len(urls) > 0 {
		container = append(container, b.CreateFlexImagesMessage(urls)...)
	}

	outerContainer := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: container,
	}

	return linebot.NewFlexMessage(MessageBlogUpdate, outerContainer)
}

func (b *Linebot) CreateFlexTextMessage(message string) *linebot.BubbleContainer {
	container := MicroBubbleContainer
	container.Body = &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeVertical,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:   linebot.FlexComponentTypeText,
				Text:   message,
				Wrap:   true,
				Weight: linebot.FlexTextWeightTypeBold,
			},
		},
	}
	return &container
}

func (b *Linebot) CreateFlexImagesMessage(urls []string) []*linebot.BubbleContainer {
	contents := []*linebot.BubbleContainer{}
	for _, url := range urls {
		content := MicroBubbleContainer

		content.Body = &linebot.BoxComponent{
			Type:       linebot.FlexComponentTypeImage,
			Layout:     linebot.FlexBoxLayoutTypeVertical,
			PaddingAll: "0px",
			Contents: []linebot.FlexComponent{
				&linebot.ImageComponent{
					Type:        linebot.FlexComponentTypeImage,
					URL:         url,
					Size:        linebot.FlexImageSizeTypeFull,
					AspectRatio: linebot.FlexImageAspectRatioType3to4,
					AspectMode:  linebot.FlexImageAspectModeTypeCover,
				},
				&linebot.BoxComponent{
					Type:            linebot.FlexComponentTypeBox,
					Layout:          linebot.FlexBoxLayoutTypeVertical,
					Position:        linebot.FlexComponentPositionTypeAbsolute,
					BackgroundColor: "#03303Acc",
					OffsetBottom:    "0px",
					OffsetStart:     "0px",
					OffsetEnd:       "0px",
					PaddingAll:      "12px",
					Contents: []linebot.FlexComponent{
						&linebot.BoxComponent{
							Type:   linebot.FlexComponentTypeBox,
							Layout: linebot.FlexBoxLayoutTypeVertical,
							Contents: []linebot.FlexComponent{
								&linebot.ButtonComponent{
									Type: linebot.FlexComponentTypeButton,
									Action: &linebot.URIAction{
										URI:   url,
										Label: "View detail",
									},
								},
							},
						},
					},
				},
			},
		}
		contents = append(contents, &content)
	}

	return contents
}

func (b *Linebot) PushMessages(ctx context.Context, to []string, messages ...linebot.SendingMessage) error {
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
