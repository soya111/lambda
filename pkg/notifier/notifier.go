package notifier

import (
	"context"
	"fmt"
	"notify/pkg/blog"
	"notify/pkg/infrastructure/line"
	"notify/pkg/model"
)

func Excute(ctx context.Context, s blog.Scraper, client *line.Linebot, subscriber model.SubscriberRepository, diary model.DiaryRepository) error {
	latestDiaries, err := getLatestDiaries(s, subscriber, diary)
	if err != nil {
		return fmt.Errorf("error getting latest diaries: %v", err)
	}

	if err := notifyAllSubscribers(ctx, client, subscriber, latestDiaries); err != nil {
		return err
	}

	return nil
}

func getLatestDiaries(s blog.Scraper, subscriber model.SubscriberRepository, diary model.DiaryRepository) ([]*blog.ScrapedDiary, error) {
	latestDiaries, err := s.ScrapeLatestDiaries()
	if err != nil {
		return nil, fmt.Errorf("error scraping latest diaries: %v", err)
	}

	res := []*blog.ScrapedDiary{}
	for _, d := range latestDiaries {
		_, err := diary.GetDiary(d.MemberName, d.Id)
		if err != nil {
			// Check if the error is a "not found" error.
			if err == model.ErrDiaryNotFound {
				// The item is not in the database, so it's a new diary.
				res = append(res, d)
				continue
			}
			// Some other error occurred.
			return nil, err
		}
	}

	return res, nil
}

func notifyAllSubscribers(ctx context.Context, client *line.Linebot, subscriber model.SubscriberRepository, diaries []*blog.ScrapedDiary) error {
	for _, d := range diaries {
		if err := notifySubscriber(ctx, client, subscriber, d); err != nil {
			return err
		}
	}
	return nil
}

func notifySubscriber(ctx context.Context, client *line.Linebot, subscriber model.SubscriberRepository, diary *blog.ScrapedDiary) error {
	to, err := subscriber.GetAllByMemberName(model.NormalizeName(diary.MemberName))
	if err != nil {
		return fmt.Errorf("error getting all by member name: %v", err)
	}

	message := line.CreateFlexMessage(diary)

	if err := client.PushMessages(ctx, to, message); err != nil {
		return fmt.Errorf("error pushing messages: %v", err)
	}

	return nil
}
