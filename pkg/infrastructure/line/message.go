package line

import (
	"fmt"
	"notify/pkg/blog"
	"notify/pkg/model"
	"notify/pkg/profile"
	"strconv"
	"unicode/utf8"

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
	sendingMessages := []linebot.SendingMessage{}
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

	message := linebot.NewFlexMessage(MessageBlogUpdate, outerContainer).WithSender(linebot.NewSender(diary.MemberName, diary.MemberIcon))
	quickReply := newQuickReplies(diary)
	message.WithQuickReplies(quickReply)

	return message
}

func createFlexTextMessage(diary *blog.ScrapedDiary) *linebot.BubbleContainer {
	container := MegaBubbleContainer

	container.Body = &linebot.BoxComponent{
		Type:       linebot.FlexComponentTypeBox,
		Layout:     linebot.FlexBoxLayoutTypeVertical,
		Height:     "400px",
		PaddingAll: "0px",
		Contents: []linebot.FlexComponent{
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Height: "70%",
				Contents: []linebot.FlexComponent{
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: []linebot.FlexComponent{
							&linebot.ImageComponent{
								Type:       linebot.FlexComponentTypeImage,
								URL:        diary.MemberIcon,
								Size:       linebot.FlexImageSizeTypeFull,
								AspectMode: linebot.FlexImageAspectModeTypeCover,
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

	if diary.IsNew() {
		// バッチのコンポーネント
		newLabel := CreateLabelComponent("NEW", "#ffffff", "#EC3D44")
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
					BackgroundColor: "#03303A",
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
									Color: "#ffffff",
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

// CreateProfileFlexMessageはプロフィールメッセージを生成
func CreateProfileFlexMessage(name string, prof *profile.Profile) linebot.SendingMessage {
	var container []*linebot.BubbleContainer
	container = append(container, createFlexProfileMessage(name, prof))

	outerContainer := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: container,
	}

	message := linebot.NewFlexMessage(name+"のプロフィール", outerContainer).WithSender(linebot.NewSender(name, prof.ImageUrl))

	return message
}

func createFlexProfileMessage(name string, prof *profile.Profile) *linebot.BubbleContainer {
	container := MegaBubbleContainer

	container.Body = &linebot.BoxComponent{
		Type:       linebot.FlexComponentTypeBox,
		Layout:     linebot.FlexBoxLayoutTypeVertical,
		Height:     "440px",
		PaddingAll: "0px",
		Contents: []linebot.FlexComponent{
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Height: "70%",
				Contents: []linebot.FlexComponent{
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: []linebot.FlexComponent{
							&linebot.ImageComponent{
								Type:       linebot.FlexComponentTypeImage,
								URL:        prof.ImageUrl,
								Size:       linebot.FlexImageSizeTypeFull,
								AspectMode: linebot.FlexImageAspectModeTypeCover,
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
										Size:   linebot.FlexTextSizeTypeSm,
										Wrap:   true,
										Text:   name,
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
									},
									&linebot.TextComponent{
										Type:   linebot.FlexComponentTypeText,
										Size:   linebot.FlexTextSizeTypeSm,
										Wrap:   true,
										Text:   fmt.Sprintf("生年月日:%s", prof.Birthday),
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
									},
									&linebot.TextComponent{
										Type:   linebot.FlexComponentTypeText,
										Size:   linebot.FlexTextSizeTypeSm,
										Wrap:   true,
										Text:   fmt.Sprintf("年齢:%s歳", prof.Age),
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
									},
									&linebot.TextComponent{
										Type:   linebot.FlexComponentTypeText,
										Size:   linebot.FlexTextSizeTypeSm,
										Wrap:   true,
										Text:   fmt.Sprintf("星座:%s", prof.Sign),
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
									},
									&linebot.TextComponent{
										Type:   linebot.FlexComponentTypeText,
										Size:   linebot.FlexTextSizeTypeSm,
										Wrap:   true,
										Text:   fmt.Sprintf("身長:%s", prof.Height),
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
									},
									&linebot.TextComponent{
										Type:   linebot.FlexComponentTypeText,
										Size:   linebot.FlexTextSizeTypeSm,
										Wrap:   true,
										Text:   fmt.Sprintf("出身地:%s", prof.Birthplace),
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
									},
									&linebot.TextComponent{
										Type:   linebot.FlexComponentTypeText,
										Size:   linebot.FlexTextSizeTypeSm,
										Wrap:   true,
										Text:   fmt.Sprintf("血液型:%s", prof.Bloodtype),
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
									},
									&linebot.ButtonComponent{
										Type:   linebot.FlexComponentTypeButton,
										Action: NewSubscribeAction(name),
										Margin: linebot.FlexComponentMarginTypeMd,
										Style:  linebot.FlexButtonStyleTypeSecondary,
										Color:  "#ffffff",
									},
								},
							},
						},
					},
				},
				PaddingAll:      "20px",
				BackgroundColor: "#464F69",
				Action: &linebot.URIAction{
					Label: "action",
					URI:   "https://www.hinatazaka46.com/s/official/artist/" + model.MemberToIdMap[name] + "?ima=0000",
				},
				Position:     linebot.FlexComponentPositionTypeAbsolute,
				OffsetBottom: "0px",
				OffsetStart:  "0px",
				OffsetEnd:    "0px",
			},
		},
	}

	generationLabelText := model.MemberToGenerationMap[name] + "期生"
	generationLabel := CreateLabelComponent(generationLabelText, "#ffffff", "#EC3D44")
	firstBox := container.Body.Contents[0].(*linebot.BoxComponent)
	firstBox.Contents = append(firstBox.Contents, generationLabel)

	return &container
}

// CreateNicknameListFlexMessageはニックネームリストを生成
func CreateNicknameListFlexMessage(name string, prof *profile.Profile) linebot.SendingMessage {
	var container []*linebot.BubbleContainer
	container = append(container, createFlexListMessage(name, prof))

	outerContainer := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: container,
	}

	message := linebot.NewFlexMessage(name+"のニックネーム", outerContainer).WithSender(linebot.NewSender(name, prof.ImageUrl))

	return message
}

func createFlexListMessage(name string, prof *profile.Profile) *linebot.BubbleContainer {
	container := MegaBubbleContainer
	components := []linebot.FlexComponent{}
	for _, nickname := range model.MemberToNicknameMap[name] {
		component := &linebot.TextComponent{
			Type:   linebot.FlexComponentTypeText,
			Size:   linebot.FlexTextSizeTypeSm,
			Wrap:   true,
			Text:   fmt.Sprintf("·%s", nickname),
			Color:  "#ffffff",
			Weight: linebot.FlexTextWeightTypeBold,
		}
		components = append(components, component)
	}

	container.Body = &linebot.BoxComponent{
		Type:       linebot.FlexComponentTypeBox,
		Layout:     linebot.FlexBoxLayoutTypeVertical,
		PaddingAll: "0px",
		Contents: []linebot.FlexComponent{
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Height: "400px",
				Contents: []linebot.FlexComponent{
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeHorizontal,
						Contents: []linebot.FlexComponent{
							&linebot.ImageComponent{
								Type:       linebot.FlexComponentTypeImage,
								URL:        prof.ImageUrl,
								Size:       linebot.FlexImageSizeTypeFull,
								AspectMode: linebot.FlexImageAspectModeTypeCover,
							},
						},
						PaddingAll: "0px",
					},
				},
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
										Size:   linebot.FlexTextSizeTypeMd,
										Wrap:   true,
										Text:   fmt.Sprintf("%sの主なニックネーム", name),
										Color:  "#ffffff",
										Weight: linebot.FlexTextWeightTypeBold,
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
											&linebot.BoxComponent{
												Type:     linebot.FlexComponentTypeBox,
												Layout:   linebot.FlexBoxLayoutTypeVertical,
												Contents: components,
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
				Position:        linebot.FlexComponentPositionTypeAbsolute,
				OffsetBottom:    "0px",
				OffsetStart:     "0px",
				OffsetEnd:       "0px",
			},
		},
	}

	generationLabelText := model.MemberToGenerationMap[name] + "期生"
	generationLabel := CreateLabelComponent(generationLabelText, "#ffffff", "#EC3D44")
	firstBox := container.Body.Contents[0].(*linebot.BoxComponent)
	firstBox.Contents = append(firstBox.Contents, generationLabel)

	return &container
}

func CreateLabelComponent(text, textcolor, backgroundcolor string) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeHorizontal,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:    linebot.FlexComponentTypeText,
				Text:    text,
				Size:    linebot.FlexTextSizeTypeXs,
				Color:   textcolor,
				Align:   linebot.FlexComponentAlignTypeCenter,
				Gravity: linebot.FlexComponentGravityTypeCenter,
			},
		},
		BackgroundColor: backgroundcolor,
		PaddingAll:      "2px",
		PaddingStart:    "4px",
		PaddingEnd:      "4px",
		Flex:            linebot.IntPtr(0),
		Position:        linebot.FlexComponentPositionTypeAbsolute,
		OffsetStart:     "18px",
		OffsetTop:       "18px",
		CornerRadius:    "100px",
		Width:           strconv.Itoa(8*utf8.RuneCountInString(text)+24) + "px",
		Height:          "25px",
	}
}

func newQuickReplies(diary *blog.ScrapedDiary) *linebot.QuickReplyItems {
	quickReplies := linebot.NewQuickReplyItems(
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("👍", "👍")),
		linebot.NewQuickReplyButton("", linebot.NewMessageAction("👎", "👎")),
	)

	if subscribeAction := NewSubscribeAction(diary.MemberName); subscribeAction != nil {
		quickReplies.Items = append(quickReplies.Items, linebot.NewQuickReplyButton("", subscribeAction))
	}

	if unsubscribeAction := newUnsubscribeAction(diary.MemberName); unsubscribeAction != nil {
		quickReplies.Items = append(quickReplies.Items, linebot.NewQuickReplyButton("", unsubscribeAction))
	}

	return quickReplies
}
