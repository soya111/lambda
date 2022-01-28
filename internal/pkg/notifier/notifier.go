package notifier

import (
	"fmt"
	"notify/internal/pkg/blog"
	"notify/internal/pkg/line"
)

func Excute(s blog.Scraper, to []string) {
	latestDiaries := s.GetAndPostLatestDiaries()
	for _, diary := range latestDiaries {
		text := fmt.Sprintf("%s %s %s\n%s", diary.Date, diary.Title, diary.MemberName, diary.Url)
		line.PushTextMessages(to, text)
		images := s.GetImages(diary)
		if len(images) > 0 {
			line.PushFlexImagesMessage(to, images)
		}
	}
}
