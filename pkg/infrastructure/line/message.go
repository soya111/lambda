package line

import (
	"fmt"
	"notify/pkg/blog"

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
	var sendingMessages []linebot.SendingMessage = []linebot.SendingMessage{}
	for _, message := range messages {
		sendingMessages = append(sendingMessages, linebot.NewTextMessage(message))
	}
	return sendingMessages
}

// CreateFlexMessages creates flex messages.
func CreateFlexMessage(diary *blog.ScrapedDiary) linebot.SendingMessage {
	var container []*linebot.BubbleContainer
	container = append(container, createFlexTextMessage(diary, true))

	container = append(container, createFlexImagesMessage(diary.Images)...)

	outerContainer := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: container,
	}

	message := linebot.NewFlexMessage(MessageBlogUpdate, outerContainer).WithSender(linebot.NewSender(diary.MemberName, diary.MemberIcon))
	quickReply := createQuickReplies()
	message.WithQuickReplies(quickReply)

	return message
}

func createFlexTextMessage(diary *blog.ScrapedDiary, showNewLabel bool) *linebot.BubbleContainer {
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

	if showNewLabel {
		// バッチのコンポーネント
		newLabel := createNewLabelComponent()
		firstBox := container.Body.Contents[0].(*linebot.BoxComponent)
		firstBox.Contents = append(firstBox.Contents, newLabel)
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

func createNewLabelComponent() *linebot.BoxComponent {
	return &linebot.BoxComponent{
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
	}
}

func createQuickReplies() *linebot.QuickReplyItems {
	quickReplies := linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("👍", "👍")),
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("👎", "👎")),
	)

	dataString, err := NewPostbackDataString(PostbackActionRegister, nil)
	if err != nil {
		fmt.Printf("createQuickReplies: %v", err)
	} else {
		registerAction := NewPostbackAction("購読する", dataString, "購読する")
		quickReplies.Items = append(quickReplies.Items, linebot.NewQuickReplyButton("", registerAction))
	}

	return quickReplies
}
