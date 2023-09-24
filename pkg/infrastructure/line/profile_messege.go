package line

import (
	"fmt"
	"zephyr/pkg/model"
	"zephyr/pkg/profile"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// CreateProfileFlexMessageはプロフィールメッセージを生成
func CreateProfileFlexMessage(prof *profile.Profile) linebot.SendingMessage {
	content := createFlexProfileMessage(prof)

	message := linebot.NewFlexMessage(prof.Name+"のプロフィール", content).WithSender(linebot.NewSender(prof.Name, prof.ImageUrl))

	return message
}

func createFlexProfileMessage(prof *profile.Profile) *linebot.BubbleContainer {
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
										Text:   prof.Name,
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
										Action: NewSubscribeAction(prof.Name),
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
					URI:   "https://www.hinatazaka46.com/s/official/artist/" + model.MemberToIdMap[prof.Name] + "?ima=0000",
				},
				Position:     linebot.FlexComponentPositionTypeAbsolute,
				OffsetBottom: "0px",
				OffsetStart:  "0px",
				OffsetEnd:    "0px",
			},
		},
	}

	generationLabelText := model.MemberToGenerationMap[prof.Name] + "期生"
	generationLabel := CreateLabelComponent(generationLabelText, "#ffffff", "#EC3D44")
	firstBox := container.Body.Contents[0].(*linebot.BoxComponent)
	firstBox.Contents = append(firstBox.Contents, generationLabel)

	return &container
}
