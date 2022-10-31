package line

import (
	"fmt"
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type Linebot struct {
	*linebot.Client
}

// 本番用コンストラクタ
func NewLinebot() *Linebot {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return &Linebot{bot}
}

// line送信
func (b *Linebot) PushTextMessages(to []string, messages ...string) {
	for _, message := range messages {
		for _, to := range to {
			if _, err := b.Client.PushMessage(to, linebot.NewTextMessage(message)).Do(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func (b *Linebot) PushFlexImagesMessage(to []string, urls []string) {
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
		if _, err := b.Client.PushMessage(to, linebot.NewFlexMessage("新着ブログがあります", container)).Do(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (b *Linebot) ReplyTextMessages(token string, message string) error {
	if _, err := b.Client.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
		return err
	}
	return nil
}
