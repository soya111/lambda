package notifier

import (
	"context"
	"fmt"
	"notify/pkg/blog"
	"notify/pkg/database"
	"notify/pkg/infrastructure"
	"os"
	"strings"
)

func Excute(ctx context.Context, s blog.Scraper, client infrastructure.Client, subscriber database.SubscriberRepository) {
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
		client.PushTextMessages(ctx, to, text)
		images := s.GetImages(diary)
		if len(images) > 0 {
			client.PushFlexImagesMessage(ctx, to, images)
		}
	}
}
