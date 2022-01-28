package notifier

import (
	"fmt"
	"notify/internal/pkg/blog"
	"notify/internal/pkg/database"
	"notify/internal/pkg/line"
	"os"
)

func ExcuteFunction() {
	to, err := database.GetDestination()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	excute(&blog.Executor{}, to)
}

func excute(s blog.ExecutorInterface, to []string) {
	latestDiaries := s.GetLatestDiaries()
	for _, diary := range latestDiaries {
		images := diary.GetImages()
		text := fmt.Sprintf("%s %s %s\n%s", diary.Date, diary.Title, diary.MemberName, diary.Url)
		line.PushTextMessages(to, text)
		line.PushFlexImagesMessage(to, images)
	}
}
