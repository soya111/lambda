package line

import (
	"fmt"
	"notify/pkg/model"
	"notify/pkg/profile"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// CreateNicknameListFlexMessageはニックネームリストを生成
func CreateNicknameListFlexMessage(prof *profile.Profile) linebot.SendingMessage {
	content := createFlexListMessage(prof)

	message := linebot.NewFlexMessage(prof.Name+"のニックネーム", content).WithSender(linebot.NewSender(prof.Name, prof.ImageUrl))

	return message
}

func createFlexListMessage(prof *profile.Profile) *linebot.BubbleContainer {
	container := MegaBubbleContainer
	components := []linebot.FlexComponent{}
	for _, nickname := range model.MemberToNicknameMap[prof.Name] {
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
										Text:   fmt.Sprintf("%sの主なニックネーム", prof.Name),
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

	generationLabelText := model.MemberToGenerationMap[prof.Name] + "期生"
	generationLabel := CreateLabelComponent(generationLabelText, "#ffffff", "#EC3D44")
	firstBox := container.Body.Contents[0].(*linebot.BoxComponent)
	firstBox.Contents = append(firstBox.Contents, generationLabel)

	return &container
}
