package notifier

import (
	"context"
	"errors"
	"fmt"
	"zephyr/pkg/blog"
	"zephyr/pkg/infrastructure/line"
	"zephyr/pkg/logging"
	"zephyr/pkg/model"

	"go.uber.org/zap"
)

// Notifier is a struct that notifies subscribers of new diaries.
type Notifier struct {
	scraper    blog.Scraper
	client     *line.Linebot
	subscriber model.SubscriberRepository
	diary      model.DiaryRepository
}

// NewNotifier creates a new Notifier.
func NewNotifier(scraper blog.Scraper, client *line.Linebot, subscriber model.SubscriberRepository, diary model.DiaryRepository) *Notifier {
	return &Notifier{
		scraper:    scraper,
		client:     client,
		subscriber: subscriber,
		diary:      diary,
	}
}

// Execute executes the notifier.
func (n *Notifier) Execute(ctx context.Context) error {
	latestDiaries, err := n.getLatestDiaries(ctx)
	if err != nil {
		return fmt.Errorf("error getting latest diaries: %v", err)
	}

	if err := n.notifyAllSubscribers(ctx, latestDiaries); err != nil {
		return err
	}

	return nil
}

// getLatestDiaries gets the latest diaries from the scraper and stores them in the database.
func (n *Notifier) getLatestDiaries(ctx context.Context) ([]*blog.ScrapedDiary, error) {
	logger := logging.LoggerFromContext(ctx)
	logger.Info("Getting latest diaries...")
	latestDiaries, err := n.scraper.ScrapeLatestDiaries(ctx)
	if err != nil {
		return nil, fmt.Errorf("error scraping latest diaries: %v", err)
	}

	res := []*blog.ScrapedDiary{}
	for _, d := range latestDiaries {
		_, err := n.diary.Get(d.MemberName, d.Id)
		if err != nil {
			// Check if the error is a "not found" error.
			if errors.Is(err, model.ErrDiaryNotFound) {
				// The item is not in the database, so it's a new diary.
				res = append(res, d)
				logger.Info("New diary", zap.Any("diary", d))
				continue
			}
			// Some other error occurred.
			return nil, err
		}
	}

	for _, sd := range res {
		diary := blog.ConvertScrapedDiaryToDiary(sd)
		if err := n.diary.Post(diary); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (n *Notifier) notifyAllSubscribers(ctx context.Context, diaries []*blog.ScrapedDiary) error {
	for _, d := range diaries {
		if err := n.notifySubscriber(ctx, d); err != nil {
			return err
		}
	}
	return nil
}

func (n *Notifier) notifySubscriber(ctx context.Context, diary *blog.ScrapedDiary) error {
	to, err := n.subscriber.GetAllByMemberName(model.NormalizeName(diary.MemberName))
	if err != nil {
		return fmt.Errorf("error getting all by member name: %v", err)
	}

	message := line.CreateFlexMessage(diary)

	if err := n.client.PushMessages(ctx, to, message); err != nil {
		return fmt.Errorf("error pushing messages: %v", err)
	}

	return nil
}
