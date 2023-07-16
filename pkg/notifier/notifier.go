package notifier

import (
	"context"
	"fmt"
	"notify/pkg/blog"
	"notify/pkg/infrastructure/line"
	"notify/pkg/infrastructure/scrape"
	"notify/pkg/model"
	"strings"
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

		document, err := scrape.GetDocumentFromURL(diary.Url)
		if err != nil {
			return fmt.Errorf("error getting document from url: %v", err)
		}
		images := s.GetImages(document)
		memberIcon := s.GetMemberIcon(document)
		message := client.CreateFlexMessage(diary, memberIcon, images)

		err = client.PushMessages(ctx, to, message)
		if err != nil {
			return fmt.Errorf("error pushing messages: %v", err)
		}
	}

	return nil
}
