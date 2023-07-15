package notifier

import (
	"context"
	"fmt"
	"notify/pkg/blog"
	"notify/pkg/infrastructure/line"
	"notify/pkg/model"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func Excute(ctx context.Context, s blog.Scraper, client *line.Linebot, subscriber model.SubscriberRepository) {
	latestDiaries, err := s.GetLatestDiaries()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	err = s.PostDiaries(latestDiaries)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, diary := range latestDiaries {
		to, err := subscriber.GetAllByMemberName(strings.Replace(diary.MemberName, " ", "", 1))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		text := fmt.Sprintf("%s %s %s\n%s", diary.Date, diary.Title, diary.MemberName, diary.Url)

		messages := []linebot.SendingMessage{}

		messages = append(messages, client.CreateTextMessages(text)...)
		images := s.GetImages(diary)
		if len(images) > 0 {
			messages = append(messages, client.CreateFlexImagesMessage(images))
		}
		err = client.PushMessages(ctx, to, messages)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	}
}
