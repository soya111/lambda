package notifier

import (
	"context"
	"fmt"
	"notify/pkg/blog"
	"notify/pkg/infrastructure/line"
	"notify/pkg/model"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func Excute(ctx context.Context, s blog.Scraper, client *line.Linebot, subscriber model.SubscriberRepository) error {
	latestDiaries, err := s.GetLatestDiaries()
	if err != nil {
		return fmt.Errorf("error getting latest diaries: %v", err)
	}

	err = s.PostDiaries(latestDiaries)
	if err != nil {
		return fmt.Errorf("error posting diaries: %v", err)
	}

	for _, diary := range latestDiaries {
		to, err := subscriber.GetAllByMemberName(strings.Replace(diary.MemberName, " ", "", 1))
		if err != nil {
			return fmt.Errorf("error getting all by member name: %v", err)
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
			return fmt.Errorf("error pushing messages: %v", err)
		}
	}

	return nil
}
