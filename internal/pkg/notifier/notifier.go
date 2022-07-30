package notifier

import (
	"fmt"
	"notify/internal/pkg/blog"
	"notify/internal/pkg/infrastructure"
	"os"
	"strings"
)

func Excute(s blog.Scraper, client infrastructure.Client, database infrastructure.Database) {
	latestDiaries := s.GetAndPostLatestDiaries()
	for _, diary := range latestDiaries {
		to, err := database.GetDestination(strings.Replace(diary.MemberName, " ", "", 1))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		text := fmt.Sprintf("%s %s %s\n%s", diary.Date, diary.Title, diary.MemberName, diary.Url)
		client.PushTextMessages(to, text)
		images := s.GetImages(diary)
		if len(images) > 0 {
			client.PushFlexImagesMessage(to, images)
		}
	}
}
