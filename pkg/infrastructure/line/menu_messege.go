package line

import (
	"notify/pkg/model"
	"notify/pkg/profile"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// CreateMenuFlexMessageはメニュー画面を生成
func CreateMenuFlexMessage(prof *profile.Profile) linebot.SendingMessage {
	content := createFlexMenuMessage(prof)

	message := linebot.NewFlexMessage(prof.Name+"のメニュー", content).WithSender(linebot.NewSender(prof.Name, prof.ImageUrl))

	return message
}

func createFlexMenuMessage(prof *profile.Profile) *linebot.BubbleContainer {
	container := MegaBubbleContainer

	container.Body = &linebot.BoxComponent{
		Type:       linebot.FlexComponentTypeBox,
		Layout:     linebot.FlexBoxLayoutTypeVertical,
		Height:     "380px",
		PaddingAll: "0px",
		Contents: []linebot.FlexComponent{
			&linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeVertical,
				Height: "70%",
				Contents: []linebot.FlexComponent{
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeVertical,
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
				Margin: linebot.FlexComponentMarginTypeLg,
				Contents: []linebot.FlexComponent{
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeVertical,
						Contents: []linebot.FlexComponent{
							&linebot.BoxComponent{
								Type:   linebot.FlexComponentTypeBox,
								Layout: linebot.FlexBoxLayoutTypeHorizontal,
								Contents: []linebot.FlexComponent{
									&linebot.ButtonComponent{
										Type:   linebot.FlexComponentTypeButton,
										Action: NewSubscribeAction(prof.Name),
										Margin: linebot.FlexComponentMarginTypeLg,
										Style:  linebot.FlexButtonStyleTypeSecondary,
										Color:  "#ffffff",
									},
									&linebot.ButtonComponent{
										Type:   linebot.FlexComponentTypeButton,
										Action: NewBlogAction(prof.Name),
										Margin: linebot.FlexComponentMarginTypeLg,
										Style:  linebot.FlexButtonStyleTypeSecondary,
										Color:  "#ffffff",
									},
								},
							},
						},
					},
					&linebot.BoxComponent{
						Type:   linebot.FlexComponentTypeBox,
						Layout: linebot.FlexBoxLayoutTypeVertical,
						Margin: linebot.FlexComponentMarginTypeLg,
						Contents: []linebot.FlexComponent{
							&linebot.BoxComponent{
								Type:   linebot.FlexComponentTypeBox,
								Layout: linebot.FlexBoxLayoutTypeHorizontal,
								Contents: []linebot.FlexComponent{
									&linebot.ButtonComponent{
										Type:   linebot.FlexComponentTypeButton,
										Action: NewProfileAction(prof.Name),
										Margin: linebot.FlexComponentMarginTypeLg,
										Style:  linebot.FlexButtonStyleTypeSecondary,
										Color:  "#ffffff",
									},
									&linebot.ButtonComponent{
										Type:   linebot.FlexComponentTypeButton,
										Action: NewNicknameAction(prof.Name),
										Margin: linebot.FlexComponentMarginTypeLg,
										Style:  linebot.FlexButtonStyleTypeSecondary,
										Color:  "#ffffff",
									},
								},
							},
						},
					},
				},
				PaddingAll:      "15px",
				BackgroundColor: "#464F69",
				Position:        linebot.FlexComponentPositionTypeAbsolute,
				OffsetBottom:    "0px",
				OffsetStart:     "0px",
				OffsetEnd:       "0px",
			},
		},
	}

	generationLabelText := model.MemberToGenerationMap[prof.Name] + "期生"
	generationLabel := CreateLabelComponent(generationLabelText, "#ffffff", "#EC3D44")
	firstBox := container.Body.Contents[0].(*linebot.BoxComponent)
	firstBox.Contents = append(firstBox.Contents, generationLabel)

	return &container
}
