package line

import (
	"fmt"
	"notify/pkg/model"
	"notify/pkg/profile"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const iconUrl = "https://cdn.hinatazaka46.com/images/14/14d/a9bac831ed1e6a4fdd93c4271aa8a.jpg"

// CreateMenuFlexMessageは汎用のメニュー画面を生成
func CreateMenuFlexMessage() linebot.SendingMessage {
	content := createFlexMenuMessage()

	message := linebot.NewFlexMessage("日向坂メニュー", content).WithSender(linebot.NewSender("日向坂46", iconUrl))

	return message
}

func createFlexMenuMessage() *linebot.BubbleContainer {
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
								URL:        iconUrl,
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
										Action: NewSelectAction(SubscribeLabel),
										Margin: linebot.FlexComponentMarginTypeLg,
										Style:  linebot.FlexButtonStyleTypeSecondary,
										Color:  "#ffffff",
									},
									&linebot.ButtonComponent{
										Type:   linebot.FlexComponentTypeButton,
										Action: NewSelectAction(BlogLabel),
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
										Action: NewSelectAction(ProfileLabel),
										Margin: linebot.FlexComponentMarginTypeLg,
										Style:  linebot.FlexButtonStyleTypeSecondary,
										Color:  "#ffffff",
									},
									&linebot.ButtonComponent{
										Type:   linebot.FlexComponentTypeButton,
										Action: NewSelectAction(NicknameLabel),
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
	return &container
}

type Action func(string) *linebot.PostbackAction

// CreateMemberSelectFlexMessageはメンバー選択画面を生成
func CreateMemberSelectFlexMessage(action Action) linebot.SendingMessage {
	container := createFlexMemberSelectMessage(action)

	outerContainer := &linebot.CarouselContainer{
		Type:     linebot.FlexContainerTypeCarousel,
		Contents: container,
	}

	message := linebot.NewFlexMessage("メンバー選択", outerContainer).WithSender(linebot.NewSender("日向坂46", iconUrl))

	return message
}

func createFlexMemberSelectMessage(action Action) []*linebot.BubbleContainer {
	contents := []*linebot.BubbleContainer{}
	groupSize := 8

	for i := 0; i < len(model.MemberList); i += groupSize {
		content := MegaBubbleContainer
		components := []linebot.FlexComponent{}

		end := i + groupSize
		if end > len(model.MemberList) {
			end = len(model.MemberList)
		}
		group := model.MemberList[i:end]

		for _, member := range group {
			component := &linebot.BoxComponent{
				Type:   linebot.FlexComponentTypeBox,
				Layout: linebot.FlexBoxLayoutTypeHorizontal,
				Margin: linebot.FlexComponentMarginTypeMd,
				Contents: []linebot.FlexComponent{
					&linebot.TextComponent{
						Type:    linebot.FlexComponentTypeText,
						Size:    linebot.FlexTextSizeTypeMd,
						Wrap:    true,
						Gravity: linebot.FlexComponentGravityTypeCenter,
						Text:    fmt.Sprintf("·%s", member),
						Color:   "#ffffff",
						Weight:  linebot.FlexTextWeightTypeBold,
					},
					&linebot.ButtonComponent{
						Type:    linebot.FlexComponentTypeButton,
						Action:  action(member),
						Height:  linebot.FlexButtonHeightTypeSm,
						Gravity: linebot.FlexComponentGravityTypeCenter,
						Style:   linebot.FlexButtonStyleTypeSecondary,
						Color:   "#ffffff",
					},
				},
			}
			components = append(components, component)
		}

		content.Body = &linebot.BoxComponent{
			Type:       linebot.FlexComponentTypeBox,
			Layout:     linebot.FlexBoxLayoutTypeVertical,
			Height:     "410px",
			PaddingAll: "0px",
			Contents: []linebot.FlexComponent{
				&linebot.BoxComponent{
					Type:            linebot.FlexComponentTypeBox,
					Layout:          linebot.FlexBoxLayoutTypeVertical,
					Contents:        components,
					PaddingAll:      "20px",
					BackgroundColor: "#464F69",
					Position:        linebot.FlexComponentPositionTypeAbsolute,
					OffsetBottom:    "0px",
					OffsetStart:     "0px",
					OffsetEnd:       "0px",
				},
			},
		}
		contents = append(contents, &content)
	}

	return contents
}

// CreateMemberMenuFlexMessageは特定メンバーのメニュー画面を生成
func CreateMemberMenuFlexMessage(prof *profile.Profile) linebot.SendingMessage {
	content := createFlexMemberMenuMessage(prof)

	message := linebot.NewFlexMessage(prof.Name+"のメニュー", content).WithSender(linebot.NewSender(prof.Name, prof.ImageUrl))

	return message
}

func createFlexMemberMenuMessage(prof *profile.Profile) *linebot.BubbleContainer {
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

	firstBox := container.Body.Contents[0].(*linebot.BoxComponent)
	firstBox.Contents = append(firstBox.Contents, CreateGenerationLabel(prof.Name, "#ffffff", "#EC3D44"))

	return &container
}
