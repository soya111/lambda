package line

import (
	"context"
	"fmt"
	"notify/pkg/blog"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/hashicorp/go-multierror"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	MessageBlogUpdate = "新着ブログがあります"
)

var MegaBubbleContainer = linebot.BubbleContainer{
	Type: linebot.FlexContainerTypeBubble,
	Size: linebot.FlexBubbleSizeTypeMega,
}

// CreateTextMessages creates text messages.
func CreateTextMessages(messages ...string) []linebot.SendingMessage {
	var sendingMessages []linebot.SendingMessage
	for _, message := range messages {
		sendingMessages = append(sendingMessages, linebot.NewTextMessage(message))
	}
	return sendingMessages
}

// CreateFlexMessages creates flex messages.
func CreateFlexMessage(diary *blog.ScrapedDiary) linebot.SendingMessage {
	var container []*linebot.BubbleContainer
	container = append(container, createFlexTextMessage(diary))

	container = append(container, createFlexImagesMessage(diary.Images)...)

	outerContainer := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: container,
	}

	return linebot.NewFlexMessage(MessageBlogUpdate, outerContainer).WithSender(linebot.NewSender(diary.MemberName, diary.MemberIcon))
}

func createFlexTextMessage(diary *blog.ScrapedDiary) *linebot.BubbleContainer {
	container := MegaBubbleContainer

	container.Body = &linebot.BoxComponent{
		Type:       linebot.FlexComponentTypeBox,
		Layout:     linebot.FlexBoxLayoutTypeVertical,
		PaddingAll: "0px",
		Contents: []linebot.FlexComponent{
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Contents: []linebot.FlexComponent{
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: []linebot.FlexComponent{
							&linebot.ImageComponent{
								Type:        linebot.FlexComponentTypeImage,
								URL:         diary.MemberIcon,
								Size:        linebot.FlexImageSizeTypeFull,
								AspectMode:  linebot.FlexImageAspectModeTypeCover,
								AspectRatio: linebot.FlexImageAspectRatioType4to3,
							},
						},
					},
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: []linebot.FlexComponent{
							&linebot.TextComponent{
								Type:    linebot.FlexComponentTypeText,
								Text:    "NEW",
								Size:    linebot.FlexTextSizeTypeXs,
								Color:   "#ffffff",
								Align:   linebot.FlexComponentAlignTypeCenter,
								Gravity: linebot.FlexComponentGravityTypeCenter,
							},
						},
						BackgroundColor: "#EC3D44",
						PaddingAll:      "2px",
						PaddingStart:    "4px",
						PaddingEnd:      "4px",
						Flex:            linebot.IntPtr(0),
						Position:        linebot.FlexComponentPositionTypeAbsolute,
						OffsetStart:     "18px",
						OffsetTop:       "18px",
						CornerRadius:    "100px",
						Width:           "48px",
						Height:          "25px",
					},
				},
				PaddingAll: "0px",
			},
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Contents: []linebot.FlexComponent{

					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeVertical,
						Contents: []linebot.FlexComponent{
							&linebot.BoxComponent{
								Type:   linebot.FlexComponentTypeBox,
								Layout: linebot.FlexBoxLayoutTypeVertical,
								Contents: []linebot.FlexComponent{
									&linebot.TextComponent{
										Type:   linebot.FlexComponentTypeText,
										Size:   linebot.FlexTextSizeTypeXl,
										Wrap:   true,
										Text:   diary.Title,
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
									},
									&linebot.TextComponent{
										Type:  linebot.FlexComponentTypeText,
										Text:  fmt.Sprintf("%s %s", diary.MemberName, diary.Date),
										Color: "#ffffffcc",
										Size:  linebot.FlexTextSizeTypeSm,
									},
								},
								Spacing: linebot.FlexComponentSpacingTypeSm,
							},
							&linebot.BoxComponent{
								Type:   linebot.FlexComponentTypeBox,
								Layout: linebot.FlexBoxLayoutTypeVertical,
								Contents: []linebot.FlexComponent{
									&linebot.BoxComponent{
										Type:   linebot.FlexComponentTypeBox,
										Layout: linebot.FlexBoxLayoutTypeVertical,
										Contents: []linebot.FlexComponent{
											&linebot.TextComponent{
												Type:   linebot.FlexComponentTypeText,
												Size:   linebot.FlexTextSizeTypeSm,
												Wrap:   true,
												Margin: linebot.FlexComponentMarginTypeXs,
												Color:  "#ffffffde",
												Text:   diary.Lead,
											},
										},
									},
								},
								PaddingAll:      "13px",
								BackgroundColor: "#ffffff1A",
								CornerRadius:    "2px",
								Margin:          linebot.FlexComponentMarginTypeXl,
							},
						},
					},
				},
				PaddingAll:      "20px",
				BackgroundColor: "#464F69",
				Action: &linebot.URIAction{
					Label: "action",
					URI:   diary.Url,
				},
				Position:     linebot.FlexComponentPositionTypeAbsolute,
				OffsetBottom: "0px",
				OffsetStart:  "0px",
				OffsetEnd:    "0px",
			},
		},
	}

	return &container
}

func createFlexImagesMessage(urls []string) []*linebot.BubbleContainer {
	contents := []*linebot.BubbleContainer{}
	num := len(urls)
	if num > 11 {
		num = 11
	}

	for _, url := range urls[:num] {

		content := MegaBubbleContainer

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

// PushMessages push messages to line
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
