package main

import (
	"fmt"
	"os"
	"time"

	"main/blog"
	"main/database"
	"main/line"

	"github.com/aws/aws-lambda-go/lambda"
)

func excute(s blog.ExecutorInterface, to []string) {
	latestDiaries := s.GetLatestDiaries()
	for _, diary := range latestDiaries {
		images := diary.GetImages()
		text := fmt.Sprintf("%s %s %s\n%s", diary.Date, diary.Title, diary.MemberName, diary.Url)
		line.PushTextMessages(to, text)
		line.PushFlexImagesMessage(to, images)
	}
}

func excuteFunction() {
	to, err := database.GetDestination()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	excute(&blog.Executor{}, to)
}

func init() {
	// set timezone
	time.Local = time.FixedZone("JST", 9*60*60)
}

func main() {
	lambda.Start(excuteFunction)
}
